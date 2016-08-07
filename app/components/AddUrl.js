import { h, Component } from 'preact'

class AddUrl extends Component {
  playUrl() {
    const url = this.state.url
    fetch('/play', {
      method: 'POST',
      headers:  {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({path: url})
    })

    return false
  }

  getInitialState() {
    return {url: ""}
  }

  handleSubmit(e) {
    e.preventDefault()
    this.playUrl()
    this.setState({url: ""})
  }

  handleChange(e) {
    this.setState({url: e.target.value});
  }

  render() {
    return (
      <div class="add-url">
        <form onSubmit={::this.handleSubmit}>
          <input type="text" value={this.state.url} onChange={::this.handleChange} />
          <input type="submit" value="Play" />
        </form>
      </div>
    )
  }
}

export default AddUrl
