import "./style.scss";
import Player from "./components/Player"

import { h, render } from 'preact'

const player = document.getElementById("player");

render(
  <Player rootUrl="/" pollInterval={200} />,
  player
);
