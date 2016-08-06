import { h, Component } from 'preact'

class Player extends Component {
  componentDidMount() {
    this.timer = setInterval(() => {
      this.updateStatusFromServer();
    }, this.props.pollInterval);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  updateStatusFromServer() {
    this.setState({ foo: "bar" })
  }

  render() {
    return (
      <div>Hello, world. {this.state.foo}</div>
    );
  }
}

export default Player
