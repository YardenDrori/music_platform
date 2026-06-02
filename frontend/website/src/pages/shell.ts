export const idForContent: string = "content";
export const shellName: string = "shell-root";

export function renderWithShell(next: (renderIn: string) => void) {
  if (document.getElementById(shellName)) {
    next(idForContent);
    return;
  }

  document.getElementById("app")!.innerHTML = `
  <div id="${shellName}">

    <div id="${idForContent}"></div>
  </div>
  `;

  next(idForContent);
}
