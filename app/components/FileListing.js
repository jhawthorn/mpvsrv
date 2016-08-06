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

  open(e) {
    const path = this.state.path
    const fullpath = path == "/" ? e.name : path + "/" + e.name
    if (e.is_dir) {
      this.visitPath(fullpath)
    } else {
      alert({play: fullpath})
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
          <a onClick={this.openParent.bind(this)} href="#">..</a>
        </div>
        <div>
          {this.state.entries.map((e) => {
            return <div><a onClick={this.open.bind(this, e)} href="#">{e.name}</a></div>
          })}
        </div>
      </pre>
    )
  }
}

export default FileListing
