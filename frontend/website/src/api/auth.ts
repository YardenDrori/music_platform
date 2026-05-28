import { setAccessToken } from "../state";
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

export async function renewAccessToken(): Promise<void> {
  const resp = await fetch("/api/token", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });

  try {
    const body = await validateResponse(resp);
    if (!body?.accessToken) {
      throw new Error("Didn't recieve access token from server");
    }
    setAccessToken(body.accessToken);
  } catch (e) {
    throw new Error("renewing access token", { cause: e });
  }
}
