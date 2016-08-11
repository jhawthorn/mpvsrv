import EventEmitter from "ev-emitter"

var emitter = new EventEmitter();

function handleResponse(status) {
  emitter.emitEvent('update', [status])
  return status;
}

function request(url, options) {
  fetch(url, options).
    then((response) => response.json()).
    then(handleResponse)
}

export function subscribe(callback) {
  emitter.on('update', callback)
}

export function unsubscribe(callback) {
  emitter.off('update', callback)
}

export function getStatus() {
  return request('/status')
}

export function stop() {
  return request('/stop', { method: 'POST' })
}

export function toggle() {
  return request('/toggle', { method: 'POST' })
}

export function seekAbsolute(seconds) {
  return request('/seek', {
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
