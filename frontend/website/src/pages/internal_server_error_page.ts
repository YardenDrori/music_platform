export function renderInternalServerErrorPage() {
  const existing = document.getElementById("err500-style");
  if (existing) existing.remove();

  const style = document.createElement("style");
  style.id = "err500-style";
  style.textContent = `
    .err500 {
      position: fixed;
      inset: 0;
      background: #080608;
      background-image: radial-gradient(ellipse at center, #150008 0%, #080608 65%);
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      font-family: 'Courier New', monospace;
      overflow: hidden;
      user-select: none;
    }

    .err500-scanlines {
      position: absolute;
      inset: 0;
      background: repeating-linear-gradient(
        0deg,
        transparent,
        transparent 2px,
        rgba(0, 0, 0, 0.12) 2px,
        rgba(0, 0, 0, 0.12) 4px
      );
      pointer-events: none;
      z-index: 20;
    }

    .err500-vignette {
      position: absolute;
      inset: 0;
      background: radial-gradient(ellipse at center, transparent 35%, rgba(0, 0, 0, 0.75) 100%);
      pointer-events: none;
      z-index: 19;
    }

    .err500-code {
      position: relative;
      font-size: clamp(7rem, 20vw, 17rem);
      font-weight: 900;
      line-height: 1;
      letter-spacing: -0.02em;
      color: #ff4500;
      text-shadow:
        0 0 40px rgba(255, 69, 0, 0.6),
        0 0 80px rgba(255, 69, 0, 0.3),
        0 0 120px rgba(255, 40, 0, 0.15);
      animation: glitch-main 6s infinite;
    }

    .err500-code::before {
      content: '500';
      position: absolute;
      top: 0; left: 0; right: 0; bottom: 0;
      color: #ff1500;
      clip-path: inset(100% 0 0 0);
      animation: glitch-layer-a 6s infinite;
    }

    .err500-code::after {
      content: '500';
      position: absolute;
      top: 0; left: 0; right: 0; bottom: 0;
      color: #ffcc00;
      clip-path: inset(100% 0 0 0);
      animation: glitch-layer-b 6s infinite;
    }

    @keyframes glitch-main {
      0%, 86%, 93%, 100% { transform: translate(0, 0); }
      87% { transform: translate(-3px,  1px); }
      88% { transform: translate( 3px, -1px); }
      89% { transform: translate(-1px,  2px); }
      90% { transform: translate( 2px, -2px); }
      91% { transform: translate(-2px,  0);   }
      92% { transform: translate(0, 0);        }
    }

    @keyframes glitch-layer-a {
      0%, 86%, 93%, 100% { clip-path: inset(100% 0 0 0); transform: translate(0); }
      87% { clip-path: inset(15% 0 60% 0); transform: translate( 5px, 0); }
      88% { clip-path: inset(50% 0 25% 0); transform: translate(-5px, 0); }
      89% { clip-path: inset(70% 0 10% 0); transform: translate( 3px, 0); }
      90% { clip-path: inset(30% 0 45% 0); transform: translate(-3px, 0); }
      92% { clip-path: inset(100% 0 0 0);  transform: translate(0);       }
    }

    @keyframes glitch-layer-b {
      0%, 87%, 94%, 100% { clip-path: inset(100% 0 0 0); transform: translate(0); }
      88% { clip-path: inset(40% 0 35% 0); transform: translate(-4px, 0); }
      89% { clip-path: inset( 5% 0 75% 0); transform: translate( 4px, 0); }
      90% { clip-path: inset(60% 0 20% 0); transform: translate(-2px, 0); }
      91% { clip-path: inset(25% 0 50% 0); transform: translate( 2px, 0); }
      93% { clip-path: inset(100% 0 0 0);  transform: translate(0);       }
    }

    .err500-subtitle {
      font-size: clamp(0.65rem, 1.8vw, 1rem);
      letter-spacing: 0.45em;
      text-transform: uppercase;
      color: #cc6633;
      margin-top: -0.3em;
      margin-bottom: 1.8em;
      animation: txt-flicker 5s infinite;
    }

    @keyframes txt-flicker {
      0%, 92%, 97%, 100% { opacity: 1;   }
      93%                { opacity: 0.3; }
      94%                { opacity: 1;   }
      95%                { opacity: 0.5; }
      96%                { opacity: 1;   }
    }

    .err500-msg {
      font-size: clamp(0.7rem, 1.4vw, 0.9rem);
      color: #996644;
      text-align: center;
      line-height: 1.7;
      margin-bottom: 2.5em;
      max-width: 380px;
      opacity: 0.85;
    }

    .err500-progress {
      width: min(360px, 75vw);
      height: 3px;
      background: rgba(255, 80, 30, 0.12);
      border-radius: 2px;
      overflow: hidden;
    }

    .err500-bar {
      height: 100%;
      width: 100%;
      background: linear-gradient(90deg, #cc1100, #ff4500, #ff9933, #ffcc00);
      transform-origin: left center;
      animation: drain 5s linear forwards;
      box-shadow: 0 0 10px rgba(255, 100, 30, 0.9);
    }

    @keyframes drain {
      from { transform: scaleX(1); }
      to   { transform: scaleX(0); }
    }

    .err500-back {
      margin-top: 0.9em;
      font-size: 0.7rem;
      letter-spacing: 0.2em;
      text-transform: uppercase;
      color: #664422;
    }

    .ember {
      position: absolute;
      width: var(--size, 3px);
      height: var(--size, 3px);
      border-radius: 50%;
      background: #ff6b35;
      box-shadow: 0 0 4px #ff6b35, 0 0 10px rgba(255, 107, 53, 0.6);
      left: var(--x, 50%);
      bottom: -4px;
      animation: ember-rise var(--dur, 3.5s) var(--delay, 0s) infinite ease-in-out;
      pointer-events: none;
    }

    @keyframes ember-rise {
      0%   { transform: translate(0, 0)                              scale(1);   opacity: 0.9; }
      40%  { transform: translate(var(--drift, 15px), -35vh)         scale(0.75); opacity: 0.7; }
      100% { transform: translate(calc(var(--drift, 15px) * 2), -85vh) scale(0.1); opacity: 0;   }
    }
  `;
  document.head.appendChild(style);

  const embers = [
    { x: 8, size: 2, dur: 3.2, delay: -1.1, drift: 12 },
    { x: 18, size: 3, dur: 4.1, delay: -0.3, drift: -18 },
    { x: 31, size: 2, dur: 2.9, delay: -2.7, drift: 25 },
    { x: 42, size: 4, dur: 3.8, delay: -0.8, drift: -10 },
    { x: 53, size: 2, dur: 4.5, delay: -3.2, drift: 30 },
    { x: 64, size: 3, dur: 3.1, delay: -1.5, drift: -22 },
    { x: 76, size: 2, dur: 4.8, delay: -2.0, drift: 18 },
    { x: 88, size: 3, dur: 3.6, delay: -0.6, drift: -14 },
  ]
    .map(
      ({ x, size, dur, delay, drift }) =>
        `<div class="ember" style="--x:${x}%;--size:${size}px;--dur:${dur}s;--delay:${delay}s;--drift:${drift}px"></div>`,
    )
    .join("");

  document.getElementById("app")!.innerHTML = `
    <div class="err500">
      <div class="err500-scanlines"></div>
      <div class="err500-vignette"></div>
      ${embers}
      <div class="err500-code">500</div>
      <div class="err500-subtitle">Server Meltdown</div>
      <div class="err500-msg">
        We're having technical difficulties.<br>
        Sorry for the inconvenience.
      </div>
      <div class="err500-progress">
        <div class="err500-bar"></div>
      </div>
      <div class="err500-back">escaping in 5s…</div>
    </div>
  `;

  let secs = 4;
  const backEl = document.querySelector<HTMLElement>(".err500-back");
  const tick = setInterval(() => {
    if (backEl)
      backEl.textContent = secs > 0 ? `Retrying in ${secs}s…` : "going back…";
    secs--;
    if (secs < 0) clearInterval(tick);
  }, 1000);

  new Promise((r) => setTimeout(r, 5000)).then(() => {
    history.back();
  });
}
