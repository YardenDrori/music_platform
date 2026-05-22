import type { AuthResponse, RegisterRequest } from "../types/auth";

export async function register(req: RegisterRequest): Promise<AuthResponse> {
  const response = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  if (response.status >= 500) {
    throw new Error("internal server error");
  }

  const resp = await response.json();

  if (!response.ok) {
    throw new Error(resp.error);
  }

  return resp;
}
