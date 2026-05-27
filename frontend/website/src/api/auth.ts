import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
} from "../types/auth";
import { validateResponse } from "../utils";

export async function register(req: RegisterRequest): Promise<AuthResponse> {
  const response = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  return validateResponse(response);
}

export async function login(req: LoginRequest): Promise<AuthResponse> {
  const response = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  return validateResponse(response);
}
