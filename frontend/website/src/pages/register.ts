import { register } from "../api/auth";
import { siteName } from "../constants";
import { route } from "../router";
import { setAccessToken } from "../state";
import type { AuthResponse, RegisterRequest } from "../types/auth";
import { verifyValidEmail } from "../utils";

export function renderRegister(): void {
  document.querySelector("#app")!.innerHTML = `
  <div class="page-centered">
    <div class="register">
    <div class="login__logo"></div>
    <h1 class="register__heading">Create Your ${siteName} Account</h1>

    <form id="register-form" class="register__form">
        <div class="register__form-text-inputs">
          <div class="register__related-input-groups-small">
            <div class="register__input-group-small">
              <input type="text" id="firstname" name="firstname" class="register__form-textbox" placeholder=" " />
              <label class="register__form-textbox-label">First name</label>
            </div>
            <div class="register__input-group-small">
              <input type="text" id="lastname" name="lastname" class="register__form-textbox" placeholder=" " />
              <label class="register__form-textbox-label">Last name</label>
            </div>
          </div>

          <div class="register__related-input-groups">
            <div class="register__input-group">
              <input type="text" id="email" name="email" class="register__form-textbox" placeholder=" " />
              <label class="register__form-textbox-label">Email</label>
            </div>
            <div class="register__input-group">
              <input type="text" id="username" name="username" class="register__form-textbox" placeholder=" " />
              <label class="register__form-textbox-label">Username</label>
            </div>
          </div>

          <div class="register__related-input-groups">
            <div class="register__input-group">
              <input type="password" id="password" name="password" class="register__form-textbox" placeholder=" " />
              <label class="register__form-textbox-label">Password</label>
            </div>
            <div class="register__input-group">
              <input type="password" id="confirm-password" name="confirm-password" class="register__form-textbox" placeholder=" " />
              <label class="register__form-textbox-label">Confirm password</label>
            </div>
          </div>
        </div>

        <div class="register__redirect-message">
          <p class="register__login-redirect-text">Already have an account? </p>
          <button type="button" id="login-button" class="register__login-redirect">Login</button>
        </div>

        <div id="form-message-div" class="register__message-container"></div>

        <button type="submit" class="register__submit-button">Continue</button>
      </form>
    </div>
  </div>
`;

  let submitButtonBlocked: boolean = false;

  const updateButton = () => {
    (document.querySelector("[type='submit']") as HTMLButtonElement).disabled =
      !(document.getElementById("firstname") as HTMLInputElement).value ||
      !(document.getElementById("lastname") as HTMLInputElement).value ||
      !(document.getElementById("email") as HTMLInputElement).value ||
      !(document.getElementById("username") as HTMLInputElement).value ||
      !(document.getElementById("password") as HTMLInputElement).value ||
      !(document.getElementById("confirm-password") as HTMLInputElement)
        .value ||
      submitButtonBlocked;

    if (submitButtonBlocked) {
      submitButtonBlocked = false;
      return;
    }
  };
  updateButton();

  document.querySelectorAll(".register__form-textbox").forEach((input) => {
    input.addEventListener("input", updateButton);
  });

  const showErrorOnForm = async (message: string) => {
    submitButtonBlocked = true;
    updateButton();
    const submitButton = document.querySelector("[type='submit']")!;
    submitButton.classList.remove("register__submit-button--loading");

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

  document
    .querySelector("#register-form")!
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      document
        .querySelector(".register__submit-button")!
        .classList.add("register__submit-button--loading");

      const formData = new FormData(
        document.querySelector("#register-form") as HTMLFormElement,
      );

      const req: RegisterRequest = {
        email: formData.get("email") as string,
        userName: formData.get("username") as string,
        firstName: formData.get("firstname") as string,
        lastName: formData.get("lastname") as string,
        password: formData.get("password") as string,
      };

      if (!verifyValidEmail(req.email)) {
        showErrorOnForm("Invalid email address");
        return;
      }

      if (req.password !== formData.get("confirm-password")) {
        showErrorOnForm("Passwords don't match, dummy");
        return;
      }

      let resp: AuthResponse;
      try {
        resp = await register(req);
        setAccessToken(resp.accessToken);
        console.log("registration successful");
        window.history.pushState({}, "", "/");
        route();
      } catch (e) {
        showErrorOnForm((e as Error).message);
      }
    });

  document.querySelector("#login-button")!.addEventListener("click", (e) => {
    e.preventDefault();
    window.history.pushState({}, "", "/login");
    route();
  });
}
