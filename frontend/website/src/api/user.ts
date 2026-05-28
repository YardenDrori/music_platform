import { authenticatedFetch } from "../utils";
import type { User } from "../types/user";

export async function getMe(): Promise<User> {
  return authenticatedFetch("/api/me", {
    method: "GET",
  });
}

export async function patchMe(modifiedUser: User) {
  return authenticatedFetch("/api/me", {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(modifiedUser),
  });
}

export async function deleteMe() {
  return authenticatedFetch("/api/me", {
    method: "DELETE",
  });
}
