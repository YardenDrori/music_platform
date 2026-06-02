export const idForContent: string = "content";
export const shellName: string = "shell-root";

export function renderWithShell(next: (renderIn: string) => void) {
  if (document.getElementById(shellName)) {
    next(idForContent);
    return;
  }

  document.getElementById("app")!.innerHTML = `
  <div id="${shellName}" class="shell-root">
    <div class="shell__top-shell-wrapper">
      <div class="shell__top-nav-bar">
        <button type="button" class="shell__home-button"></button>
        <button type="button" class="shell__hot-button"></button>
        <button type="button" class="shell__generate-button"></button>
        <button type="button" class="shell__search-button"></button>
      </div>
      <div class="shell__account-pic-placeholder"></div>
    </div>

    <div id="${idForContent}" class="shell__content"></div>

    <div class="shell__bottom-bar">
      <div class="shell__current-playing-song">
        <div class="shell__now-playing-album-pic-placeholder"></div>
        <div class="shell__now-playing-title-and-runtime">
          <p class="shell__now-playing-song-name">Never Gonna Give You Up</p>
          <p class="shell__now-playing-artist-name">Rick Astly ft Your Mom</p>
          <div class="shell__now-playing-runtime-section">
            <p class="shell__current-playing-runtime-current">69:69</p>
            <div class="shell__current-playing-runtime-progress-bar-current"></div>
            <div class="shell__current-playing-runtime-progress-bar-total"></div>
            <p class="shell__current-playing-runtime-total">420:420</p>
          </div>
        </div>
      </div>

      <div class="shell__media-controls">
        <button type="button" class="shell__shuffle-button"></button>
        <button type="button" class="shell__prev-song-button"></button>
        <button type="button" class="shell__pause-play-button"></button>
        <button type="button" class="shell__next-song-button"></button>
        <button type="button" class="shell__loop-button"></button>
      </div>

      <div class="shell__misc-buttons-group">
        <div class="shell__misc-buttons-subgroup">
          <button type="button" class="shell__like-button"></button>
          <button type="button" class="shell__dislike-button"></button>
        </div>
        <div class="shell__misc-buttons-subgroup">
          <button type="button" class="shell__placeholder-button"></button>
          <button type="button" class="shell__volume-button"></button>
        </div>
        <div class="shell__misc-buttons-subgroup">
          <button type="button" class="shell__subtitiles-button"></button>
          <button type="button" class="shell__queue-button"></button>
        </div>
      </div>

    </div>
  </div>
  `;

  next(idForContent);
}
