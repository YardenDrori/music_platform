import { renderDashboard } from "./pages/dashboard";
import { renderLogin } from "./pages/login";
import { renderRegister } from "./pages/register";
import { getAccessToken } from "./state";

function needLogin(render: () => void) {
  if (!getAccessToken()) {
    history.pushState({}, "", "/login");
  } else {
    render();
  }
  route();
}

export function route() {
  const path = window.location.pathname;
  if (path === "/login") {
    renderLogin();
  } else if (path === "/register") {
    renderRegister();
  } else {
    //default route
    needLogin(() => renderDashboard);
  }
}
