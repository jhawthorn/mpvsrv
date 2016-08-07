import { h, Component } from 'preact'
import path from 'path'

class FileListing extends Component {
  constructor() {
    super()

    this.setState({
      path: "/",
      entries: [],
    })

    this.visitPath("/")
  }

  visitPath(path) {
    fetch('/browse/'+path)
      .then((response) => response.json())
      .then((json) => this.setState({
        path: path,
        entries: json
      }))
  }

  playPath(path) {
    fetch('/play', {
      method: 'POST',
      headers:  {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({path: path})
    })
  }

  open(e) {
    const path = this.state.path
    const fullpath = path == "/" ? e.name : path + "/" + e.name
    if (e.is_dir) {
      this.visitPath(fullpath)
    } else {
      this.playPath(fullpath)
    }
  }

  openParent() {
    const oldPath = this.state.path
    const newPath = oldPath.split("/").slice(0, -1).join("/")
    if (!newPath) {
      this.visitPath("/")
    } else {
      this.visitPath(newPath)
    }
  }

  render() {
    return (
      <pre>
        <h2>{this.state.path}</h2>
        <div>
          <a onClick={::this.openParent} href="#">..</a>
        </div>
        <div>
          {this.state.entries.map((e) => {
            return <div><a onClick={() => this.open(e)} href="#">{e.name}</a></div>
          })}
        </div>
      </pre>
    )
  }
}

export default FileListing
