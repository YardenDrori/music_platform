import { route } from "./router";

route();

window.addEventListener("popstate", route);
