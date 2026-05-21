import type { AuthResponse, RegisterRequest } from "../types/auth";

export async function register(req: RegisterRequest): Promise<AuthResponse> {
  const respones = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  const resp = await respones.json();

  if (!respones.ok) {
    throw new Error(resp.error);
  }

  return resp;
}
