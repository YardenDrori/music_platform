import { renderLogin } from "./pages/login";
import { renderRegister } from "./pages/register";

export function route() {
  const path = window.location.pathname;
  if (path === "/login") {
    renderLogin();
  } else if (path === "/register") {
    renderRegister();
  } else {
    renderLogin();
  }
}

route();

window.addEventListener("popstate", route);
