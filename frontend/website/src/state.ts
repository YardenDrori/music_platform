import { nextTrackIconA, prevTrackIconA } from "./icons";

let accessToken: string | null = null;
export let prevTrackStyle: string = prevTrackIconA;
export let nextTrackStyle: string = nextTrackIconA;

export function setAccessToken(token: string) {
  accessToken = token;
}

export function getAccessToken(): string | null {
  return accessToken;
}
