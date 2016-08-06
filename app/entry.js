import "./style.scss";
import Player from "./components/Player"

import { h, render } from 'preact'

render(
  <Player rootUrl="/" pollInterval={1000} />,
  document.body
);
