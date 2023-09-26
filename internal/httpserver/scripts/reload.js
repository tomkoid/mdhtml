var loc = window.location;
var uri = 'ws:';

if (loc.protocol === 'https:') {
  uri = 'wss:';
}
uri += '//' + loc.host;
uri += loc.pathname + 'ws';

const socket = new WebSocket(uri);

console.log('[websocket] connecting...')

socket.onopen = () => {
  console.log('[websocket] connected');
  socket.send('connected');

}

socket.onmessage = async (event) => {
  if (event.data === 'hello') {
    console.log('[websocket] server responded!');
    return
  }
  console.log(`[websocket] message: ${event.data}`);
  if (event.data === 'reload') {
    const content = await fetch('/');
    const html = await content.text();

    const parser = new DOMParser();
    const doc = parser.parseFromString(html, 'text/html');

    // use the new body instead of the old one
    document.body = doc.body;
  }
}

socket.onclose = () => {
  console.log('[websocket] disconnected');
}

socket.onerror = (error) => {
  console.log(`[websocket] error: ${error}`);
  console.log(error)
}
