import { getAccessToken } from "../state";
import { validateEmptyResponse, validateResponse } from "../utils";
import type { User } from "../types/user";

export async function getMe(): Promise<User> {
  const access = getAccessToken();

  const resp = await fetch("/api/me", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${access}`,
    },
  });

  return await validateResponse(resp);
}

export async function patchMe(modifiedUser: User) {
  const access = getAccessToken();

  const resp = await fetch("/api/me", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${access}`,
    },
    body: JSON.stringify(modifiedUser),
  });

  return await validateEmptyResponse(resp);
}

export async function deleteMe() {
  const access = getAccessToken();

  const resp = await fetch("/api/me", {
    method: "DELETE",
    headers: { Authirzation: `Bearer ${access}` },
  });

  return await validateEmptyResponse(resp);
}
