import { renewAccessToken } from "./api/auth";
import { renderDashboard } from "./pages/dashboard";
import { renderLogin } from "./pages/login";
import { renderRegister } from "./pages/register";
import { renderWithShell } from "./pages/shell";
import { getAccessToken } from "./state";

function needLogin(render: () => void) {
  if (!getAccessToken()) {
    console.log(
      "entering an auth endpoint with now access token, attempting to renew",
    );
    renewAccessToken().then(() => {
      if (!getAccessToken()) {
        console.log("failed to renew access token, directing to login");
        history.pushState({}, "", "/login");
        route();
        return;
      }
      render();
    });
  } else {
    render();
  }
}

export function route() {
  const path = window.location.pathname;
  if (path === "/login") {
    renderLogin();
  } else if (path === "/register") {
    renderRegister();
  } else {
    //default route
    needLogin(() => renderWithShell(renderDashboard));
  }
}
