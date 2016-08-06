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

  fetch() {
    fetch.apply(null, arguments)
      .then((response) => response.json())
      .then((json) => this.setState(json))
      .catch((ex) => console.log('parsing failed', ex))
  }

  render() {
    const time = this.state.time || {};
    return (
      <div>
        <h2>{this.state.title}</h2>
        <div>{time.current}/{time.total}</div>
        <div>
          <progress value={time.current} max={time.total}>{time.percent}%</progress>
        </div>
        <button onClick={this.toggle.bind(this)}>{this.state.paused ? "⏵ play" : "⏸ pause" }</button>
      </div>
    );
  }
}

export default Player
