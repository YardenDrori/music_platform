import { renewAccessToken } from "./api/auth";
import { renderDashboard } from "./pages/dashboard";
import { renderInternalServerErrorPage } from "./pages/internal_server_error_page";
import { renderLogin } from "./pages/login";
import { renderRegister } from "./pages/register";
import { renderWithShell } from "./pages/shell";
import { getAccessToken } from "./state";
import { InternalError } from "./types/errors";

function needLogin(render: () => void) {
  if (!getAccessToken()) {
    console.log(
      "entering an auth endpoint with now access token, attempting to renew",
    );
    renewAccessToken()
      .then(() => {
        if (!getAccessToken()) {
          console.log("failed to renew access token, directing to login");
          history.pushState({}, "", "/login");
          route();
          return;
        }
        render();
      })
      .catch((e) => {
        if (e instanceof InternalError) return;
        console.log("failed to renew access token, directing to login");
        history.pushState({}, "", "/login");
        route();
        return;
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
  } else if (path === "/internal-error") {
    renderInternalServerErrorPage();
  } else {
    //default route
    needLogin(() => renderWithShell(renderDashboard));
  }
}
