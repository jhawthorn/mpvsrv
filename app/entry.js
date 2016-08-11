import "./style.scss";
import Player from "components/Player"
import FileListing from "components/FileListing"
import AddUrl from "components/AddUrl"
import {addKeybindings} from "keybindings"

import { h, render } from 'preact'

addKeybindings(document);

render(
  (
    <div>
      <Player rootUrl="/" pollInterval={1000} />
      <hr />
      <AddUrl />
      <hr />
      <FileListing />
    </div>
  ),
  document.body
);
