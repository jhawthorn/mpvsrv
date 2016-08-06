import { h, Component } from 'preact'

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
    fetch.apply(null, arguments)
      .then((response) => response.json())
      .then((json) => this.setState(json))
      .catch((ex) => console.log('parsing failed', ex))
  }

  render() {
    const state = this.state;
    const time = this.state.time || {};
    return (
      <div>
        <h2>{state.title}</h2>
        <div>{time.current}/{time.total}</div>
        <div>
          <progress value={time.current} max={time.total}>{time.percent}%</progress>
        </div>
        <button disabled={state.idle} onClick={this.toggle.bind(this)}>{this.state.paused ? "⏵ Play" : "⏸ Pause" }</button>
        <button disabled={state.idle} onClick={this.stop.bind(this)}>⏹ Stop</button>
      </div>
    );
  }
}

export default Player
