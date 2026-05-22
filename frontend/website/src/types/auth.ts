import type { User } from "./user";

export type RegisterRequest = {
  email: string;
  userName: string;
  firstName: string;
  lastName: string;
  password: string;
};

export type LoginRequest = {
  identifier: string;
  password: string;
};

export type AuthResponse = {
  user: User;
  accessToken: string;
};
