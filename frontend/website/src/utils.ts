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
    throw new Error("internal server error");
  }

  const respJson = await resp.json();

  if (!resp.ok) {
    throw new Error(respJson.error);
  }

  return respJson;
}
