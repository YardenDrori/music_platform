import { register } from "../api/auth";
import { setAccessToken } from "../state";
import type { AuthResponse, RegisterRequest } from "../types/auth";

export function renderRegister(): void {
  document.querySelector("#app")!.innerHTML = `
<form id="register-form">
  <h1>Create an account</h1>

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

      if (
        (formData.get("password") as string) !==
        formData.get("confirm-password")
      ) {
        alert("passwords do not match");
        return;
      }

      const req: RegisterRequest = {
        email: formData.get("email") as string,
        userName: formData.get("username") as string,
        firstName: formData.get("firstname") as string,
        lastName: formData.get("lastname") as string,
        password: formData.get("password") as string,
      };

      let resp: AuthResponse;
      try {
        resp = await register(req);
        setAccessToken(resp.accessToken);
        console.log("registration successful");
      } catch (e) {
        console.log("registration failed: " + e);
      }
    });
}
