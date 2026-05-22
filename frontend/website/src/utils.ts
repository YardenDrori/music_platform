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
