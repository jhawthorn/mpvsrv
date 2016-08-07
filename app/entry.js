import "./style.scss";
import Player from "components/Player"
import FileListing from "components/FileListing"

import { h, render } from 'preact'

render(
  (
    <div>
      <Player rootUrl="/" pollInterval={1000} />
      <hr />
      <FileListing />
    </div>
  ),
  document.body
);
