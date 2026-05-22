import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
} from "../types/auth";

async function validateResponse(resp: Response): Promise<any> {
  if (resp.status >= 500) {
    throw new Error("internal server error");
  }

  const respJson = await resp.json();

  if (!resp.ok) {
    throw new Error(respJson.error);
  }

  return respJson;
}

export async function register(req: RegisterRequest): Promise<AuthResponse> {
  const response = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  return validateResponse(response);
}

  if (!response.ok) {
    throw new Error(resp.error);
  }

  return resp;
}
