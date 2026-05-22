import { login } from "../api/auth";
import { route } from "../router";
import { setAccessToken } from "../state";
import type { AuthResponse, LoginRequest } from "../types/auth";
import { verifyValidEmail } from "../utils";

export function renderLogin(): void {
  document.querySelector("#app")!.innerHTML = `
  <div>
    <form id="login-form">
      <h1>Welcome back! Please login.</h1>

      <div>
        <label for="identifier">Email or username</label>
        <input type="text" id="identifier" name="identifier" required />
      </div>
      <div>
        <label for="password">password</label>
        <input type="password" id="password" name="password" required />
      </div>

      <p id="form-message"></p>
      <button type="submit">Login</button>
    </form>
    <a id="register-button">Register</a>
  </div>
  `;

  document
    .querySelector("#login-form")!
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      const formData = new FormData(
        document.querySelector("#login-form")! as HTMLFormElement,
      );

      const identifier = formData.get("identifier") as string;
      let req: LoginRequest;

      if (identifier.includes("@")) {
        //email
        if (!verifyValidEmail(identifier)) {
          document.querySelector("#form-message")!.textContent =
            "Invalid Email address";
          return;
        }

        req = {
          email: identifier as string,
          username: "",
          password: formData.get("password") as string,
        };
      } else {
        //username
        req = {
          email: "",
          username: identifier as string,
          password: formData.get("password") as string,
        };
      }

      let resp: AuthResponse;
      try {
        resp = await login(req);
        console.log("login successful");
      } catch (e) {
        document.getElementById("form-message")!.textContent =
          "" + (e as Error).message;
        return;
      }

      setAccessToken(resp.accessToken);
    });

  document.querySelector("#register-button")!.addEventListener("click", (e) => {
    e.preventDefault();
    window.history.pushState({}, "", "/register");
    route();
  });
}
