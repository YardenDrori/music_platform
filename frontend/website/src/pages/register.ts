import { register } from "../api/auth";
import { setAccessToken } from "../state";
import type { AuthResponse, RegisterRequest } from "../types/auth";

export function renderRegister(): void {
  document.querySelector("#app")!.innerHTML = `
<form id="register-form">
<input type="email" name="email" placeholder="Email" />
<input type="text" name="username" placeholder="Username" />
<input type="text" name="firstname" placeholder="Firstname" />
<input type="text" name="lastname" placeholder="Lastname" />
<input type="password" name="password" placeholder="Password" />
<input type="text" name="confirm-password" placeholder="Confirm Password" />
<button>Submit</button>
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
