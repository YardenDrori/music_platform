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

  const respJson = await resp.json();

  if (!resp.ok) {
    throw new Error(respJson?.error ?? "unkown error");
  }

  return respJson;
}

export async function validateEmptyResponse(resp: Response): Promise<void> {
  if (resp.status >= 500) {
    throw new Error("Internal server error");
  }
  if (!resp.ok) {
    const respJson = await resp.json();
    throw new Error(respJson?.error ?? "unkown error");
  }
}
