import { getMe } from "../api/user";

export function renderDashboard(): void {
  document.getElementById("app")!.innerHTML = `
  <p>Fetching account information, this shouldn't take long...</p>
  `;

  getMe()
    .then((user) => {
      console.log("own details fetched");
      const app = document.getElementById("app")!;
      app.innerHTML = "";

      const welcomeMsg = document.createElement("p");
      welcomeMsg.textContent = "Welcome back " + user.firstName + "!";

      app.appendChild(welcomeMsg);
    })
    .catch((err) => {
      const app = document.getElementById("app")!;
      const message = err instanceof Error ? err.message : "unknown error";
      app.textContent = "An unexpected error occurred: " + message;
    });
}
