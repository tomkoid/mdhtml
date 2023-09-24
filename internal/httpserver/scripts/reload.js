var loc = window.location;
var uri = 'ws:';

if (loc.protocol === 'https:') {
  uri = 'wss:';
}
uri += '//' + loc.host;
uri += loc.pathname + 'ws';

ws = new WebSocket(uri)

const socket = new WebSocket(uri);

socket.onopen = () => {
  console.log('[websocket] connected');
  socket.send('connected');
}

socket.onmessage = (event) => {
  if (event.data === 'hello') {
    console.log('[websocket] server responded!');
    return
  }
  console.log(`[websocket] message: ${event.data}`);
  if (event.data === 'reload') {
    socket.close()
    window.location.reload();
  }
}

socket.onclose = () => {
  console.log('[websocket] disconnected');
}

socket.onerror = (error) => {
  console.log(`[websocket] error: ${error}`);
  console.log(error)
}
