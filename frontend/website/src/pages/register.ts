import { register } from "../api/auth";
import { setAccessToken } from "../state";
import type { AuthResponse, RegisterRequest } from "../types/auth";

export function renderRegister(): void {
  document.querySelector("#app")!.innerHTML = `
<form id="register-form">
  <h1>Create an account!</h1>

  <div>
    <label for="email">Email</label>
    <input type="email" id="email" name="email" required />
  </div>

  <div>
    <label for="username">Username</label>
    <input type="text" id="username" name="username" required />
  </div>

  <div>
    <label for="firstname">First name</label>
    <input type="text" id="firstname" name="firstname" required />
  </div>

  <div>
    <label for="lastname">Last name</label>
    <input type="text" id="lastname" name="lastname" required />
  </div>

  <div>
    <label for="password">Password</label>
    <input type="password" id="password" name="password" required />
  </div>

  <div>
    <label for="confirm-password">Confirm password</label>
    <input type="password" id="confirm-password" name="confirm-password" required />
  </div>

  <p id="form-message"></p>
  <button type="submit">Register</button>
</form>
`;

  document
    .querySelector("#register-form")!
    .addEventListener("submit", async (e) => {
      e.preventDefault();

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

      if (req.password !== formData.get("confirm-password")) {
        document.querySelector("#form-message")!.textContent =
          "Passwords do not match dumbass lmao";
        return;
      }

      let resp: AuthResponse;
      try {
        resp = await register(req);
        setAccessToken(resp.accessToken);
        document.querySelector("#form-message")!.textContent =
          "Registration successful!";
      } catch (e) {
        document.querySelector("#form-message")!.textContent =
          "" + (e as Error).message;
      }
    });
}
