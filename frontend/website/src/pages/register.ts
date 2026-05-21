export function renderRegister(): void {
  document.querySelector("#app")!.innerHTML = `
<form id="register-form">
<input type="email" placeholder="Email" />
<input type="text" placeholder="Username" />
<input type="text" placeholder="Firstname" />
<input type="text" placeholder="Lastname" />
<input type="password" placeholder="Password" />
<input type="text" placeholder="Confirm Password" />
<button>Submit</button>
</form>
`;

  document.querySelector("#register-form")!.addEventListener("submit", (e) => {
    e.preventDefault();
    console.log("submitted registration request");
  });
}
