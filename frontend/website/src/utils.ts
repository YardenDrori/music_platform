import { renewAccessToken } from "./api/auth";
import { getAccessToken } from "./state";

export function verifyValidEmail(email: string): boolean {
  const identifierParts = email.split("@");

  if (identifierParts.length !== 2) {
    return false;
  }

  if (
    !identifierParts[1].includes(".") ||
    identifierParts[1].length === 0 ||
    identifierParts[0].length === 0
  ) {
    return false;
  }

  return true;
}

export async function validateResponse(resp: Response): Promise<any> {
  if (resp.status >= 500) {
    throw new Error("Internal server error");
  }

  let body: any = null;
  try {
    body = await resp.json();
  } catch (_) {}

  if (!resp.ok) {
    throw new Error(body?.error ?? "unknown error");
  }

  return body;
}

export async function authenticatedFetch(
  url: string,
  params: RequestInit,
): Promise<any> {
  const makeRequest = async (): Promise<any> => {
    const token = getAccessToken();
    if (!token) {
      throw new Error("called AuthenticatedFetch without a token");
    }

    return fetch(url, {
      ...params,
      headers: {
        ...params.headers,
        Authorization: `Bearer ${token}`,
      },
    });
  };

  let resp = await makeRequest();

  if (resp.status === 401) {
    try {
      await renewAccessToken();
      return await validateResponse(await makeRequest());
    } catch (e) {
      throw new Error(
        `making an authenticatedFetch request to ${url}, ${params}`,
        { cause: e },
      );
    }
  }

  return validateResponse(resp);
}
