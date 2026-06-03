import {
  libraryIcon,
  likeIcon,
  lyricsIcon,
  playIcon,
  queueIcon,
  repeatOffIcon,
  shuffleIcon,
  tagIcon,
} from "../icons";
import { favsTagButtonStyle, nextTrackStyle, prevTrackStyle } from "../state";

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
      <div class="shell__account-pic"></div>
    </div>

    <div id="${idForContent}" class="shell__content"></div>

    <div class="shell__bottom-bar">
      <div class="shell__bottom-bar-controls">
        <div class="shell__current-playing-song">
          <div class="shell__now-playing-album-pic"></div>
          <div class="shell__now-playing-title-and-artist">
            <button type="button" class="shell__now-playing-song-name">Never Gonna Give You Up</button>
            <button type="button" class="shell__now-playing-artist-name">Rick Astly ft Your Mom</button>
          </div>
        </div>

        <div class="shell__media-controls">
          <button type="button" class="shell__shuffle-button">${shuffleIcon}</button>
          <button type="button" class="shell__prev-song-button">${prevTrackStyle}</button>
          <button type="button" class="shell__pause-play-button">${playIcon}</button>
          <button type="button" class="shell__next-song-button">${nextTrackStyle}</button>
          <button type="button" class="shell__loop-button">${repeatOffIcon}</button>
        </div>

        <div class="shell__misc-buttons-group">
          <div class="shell__misc-buttons-subgroup">
            <button type="button" class="shell__favorites-button">${favsTagButtonStyle}</button>
            <button type="button" class="shell__rate-button">
              ${likeIcon}
            </button>
          </div>
          <div class="shell__misc-buttons-subgroup-middle">
            <button type="button" class="shell__add-tag-button">${tagIcon}</button>
            <button type="button" class="shell__add-to-library-button">${libraryIcon}</button>
          </div>
          <div class="shell__misc-buttons-subgroup">
            <button type="button" class="shell__subtitiles-button">${lyricsIcon}</button>
            <button type="button" class="shell__queue-button">${queueIcon}</button>
          </div>
        </div>
      </div>

      <div class="shell__now-playing-runtime-section">
        <div class="shell__runtime-number-wrapper-start">
          <p class="shell__current-playing-runtime-current">69:69</p>
        </div>
        <div class="shell__current-playing-runtime-progress-bar-total">
          <div class="shell__current-playing-runtime-progress-bar-current"></div>
        </div>
        <div class="shell__runtime-number-wrapper-end">
          <p class="shell__current-playing-runtime-total">420:420</p>
        </div>
      </div>

    </div>
  </div>
  `;

  next(idForContent);
}
