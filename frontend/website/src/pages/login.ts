import { login } from "../api/auth";
import { siteName } from "../constants";
import { route } from "../router";
import { setAccessToken } from "../state";
import type { AuthResponse, LoginRequest } from "../types/auth";
import { verifyValidEmail } from "../utils";

export function renderLogin(): void {
  document.querySelector("#app")!.innerHTML = `
  <div class="page-centered">
    <div class="login">
      <div class="login__logo"></div>

      <h1 class="login__heading">Sign In With ${siteName} Account</h1>

      <form id="login-form" class="login__form">
        <div class="login__form-text-inputs">
          <div class="login__input-group">
            <input type="text" id="identifier" name="identifier" class="login__form-textbox"  placeholder=" "/>
            <label class="login__form-textbox-label">Email or Username</label>
          </div>
          <div class="login__input-group">
            <input type="password" id="password" name="password" class="login__form-textbox" placeholder=" "/>
            <label class="login__form-textbox-label">Password</label>
          </div>
        </div>
        <button type="button" id="register-button" class="login__register-redirect">Create Your ${siteName} Account</button>
        <div id="form-message-div"></div>

        <button type="submit" class="login__submit-button">Continue</button>
      </form>

    </div>
  </div>
  `;

  let submitButtonBlocked: boolean = false;

  const updateButton = () => {
    (document.querySelector("[type='submit']") as HTMLButtonElement).disabled =
      !(document.getElementById("identifier") as HTMLInputElement).value ||
      !(document.getElementById("password") as HTMLInputElement).value ||
      submitButtonBlocked;

    if (submitButtonBlocked) {
      submitButtonBlocked = false;
      return;
    }
  };
  updateButton();

  const showErrorOnForm = async (message: string) => {
    submitButtonBlocked = true;
    updateButton();
    const submitButton = document.querySelector("[type='submit']")!;
    submitButton.classList.remove("login__submit-button--loading");

    const errMsg = document.getElementById("form-message");
    if (errMsg) {
      errMsg.remove();
    }

    const newMsg = document.createElement("p");
    newMsg.id = "form-message";
    newMsg.className = "login__form-message";
    newMsg.textContent = message;

    document.getElementById("form-message-div")!.appendChild(newMsg);
  };

  document.querySelectorAll(".login__form-textbox").forEach((input) => {
    input.addEventListener("input", updateButton);
  });

  document
    .querySelector("#login-form")!
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      document
        .querySelector(".login__submit-button")!
        .classList.add("login__submit-button--loading");

      const formData = new FormData(
        document.querySelector("#login-form")! as HTMLFormElement,
      );

      const identifier = formData.get("identifier") as string;
      let req: LoginRequest;

      if (identifier.includes("@")) {
        //email
        if (!verifyValidEmail(identifier)) {
          showErrorOnForm("Invalid email address");
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
        setAccessToken(resp.accessToken);
        console.log("login successful");
        window.history.pushState({}, "", "/");
        route();
      } catch (e) {
        showErrorOnForm((e as Error).message);
        return;
      }
    });

  document.querySelector("#register-button")!.addEventListener("click", (e) => {
    e.preventDefault();
    window.history.pushState({}, "", "/register");
    route();
  });
}
