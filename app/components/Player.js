import { h, Component } from 'preact'

function pad2(n) {
  if (n < 10) {
    return "0" + n
  } else {
    return "" + n
  }
}

function formatTime(t) {
  let seconds = Math.round(t)
  let minutes = seconds / 60 | 0
  let hours = minutes / 60 | 0
  seconds = pad2(seconds % 60)
  minutes = pad2(minutes % 60)
  hours = pad2(hours)
  return `${hours}:${minutes}:${seconds}`
}

class Progress extends Component {
  percent() {
    return this.current() * 100 / this.props.total;
  }

  current() {
    return this.state.preview ? this.state.preview : this.props.current
  }

  clickHandler(e) {
    const value = this.props.total * e.offsetX / this.progressBar.offsetWidth;
    this.props.onSeek(value)
  }

  mouseMoveHandler(e) {
    const value = this.props.total * e.offsetX / this.progressBar.offsetWidth;
    this.setState({preview: value})
  }

  mouseLeaveHandler(e) {
    this.setState({preview: null})
  }

  render() {
    return (
      <div class="progress">
        <div class="progress-time">
          <span class="time-current">{formatTime(this.current())}</span>
          <span class="time-of">/</span>
          <span class="time-total">{formatTime(this.props.total)}</span>
        </div>
        <div
          ref={(c) => this.progressBar = c}
          class="progress-bar"
          onClick={::this.clickHandler}
          onMouseMove={::this.mouseMoveHandler}
          onMouseLeave={::this.mouseLeaveHandler}>
          <div class="progress-bar-fill" style={`width: ${this.percent()}%`}></div>
        </div>
      </div>
    )
  }
}

class Player extends Component {
  componentDidMount() {
    this.timer = setInterval(() => {
      this.updateStatusFromServer();
    }, this.props.pollInterval);
    this.updateStatusFromServer();
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  updateStatusFromServer() {
    this.fetch('/status')
  }

  toggle() {
    this.fetch('/toggle', { method: 'POST' })
  }

  stop() {
    this.fetch('/stop', { method: 'POST' })
  }

  fetch() {
    this.handleResponse(
      fetch.apply(null, arguments).then((response) => response.json())
    )
  }

  handleResponse(dfd) {
    dfd
      .then((json) => this.setState(json))
      .catch((ex) => console.log('parsing failed', ex))
  }

  seekHandler(seconds) {
    this.fetch('/seek', {
      method: 'POST',
      headers:  {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        seconds: seconds
      })
    })
  }

  render() {
    const state = this.state;
    const time = this.state.time || {};
    return (
      <div class="player">
        <div class="player-title">{state.title}</div>
        <div class="player-progress">
          <Progress current={time.current} total={time.total} onSeek={::this.seekHandler} />
        </div>
        <div class="player-actions">
          <button class="player-action action-prev" disabled >Prev</button>
          <button class="player-action action-playpause" disabled={state.idle} onClick={::this.toggle}>{this.state.paused ? "⏵ Play" : "⏸ Pause" }</button>
          <button class="player-action action-stop" disabled={state.idle} onClick={::this.stop}>⏹ Stop</button>
          <button class="player-action action-next" disabled >Next</button>
        </div>
      </div>
    );
  }
}

export default Player
