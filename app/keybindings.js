import {toggle, stop, seekRelative} from "mpv"

const keybindings = {
  " ": () => toggle(),
  "ArrowLeft": () => seekRelative(-5),
  "ArrowRight": () => seekRelative(5)
}

function onKeydown(e) {
  const handler = keybindings[e.key];
  if (handler) {
    handler(e)
  }
}

export function addKeybindings(element) {
  element.addEventListener("keydown", onKeydown, false)
}
