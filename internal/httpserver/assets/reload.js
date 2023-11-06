let wSStatus = 3; // 0 = disconnected, 1 = connected, 2 = connecting, 3 = initial
let socket;

function connect() {
  var loc = window.location;
  var uri = 'ws:';

  if (loc.protocol === 'https:') {
    uri = 'wss:';
  }
  uri += '//' + loc.host;
  uri += loc.pathname + 'ws';

  socket = new WebSocket(uri);

  console.log('[websocket] connecting...')

  socket.onopen = () => {
    console.log('[websocket] connected');
    socket.send('connected');

  }

  socket.onmessage = async (event) => {
    if (event.data === 'hello') {
      wSStatus = 0;
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

      // reload the scripts
      const scripts = document.querySelectorAll('script');
      scripts.forEach((script) => {
        if (script.src.includes('reload.js')) {
          return;
        }

        console.log(`[DOM] reloading script: ${script.src}`);

        const newScript = document.createElement('script');
        newScript.src = script.src;
        newScript.async = false;
        document.body.appendChild(newScript);
      });
    }
  }

  socket.onclose = () => {
    wSStatus = 1;

    console.log('[websocket] connection closed');
    setTimeout(() => {
      wSStatus = 2;
      console.log('[websocket] reconnecting...');
      connect();
    }, 1000);
  }

  socket.onerror = (error) => {
    console.log(`[websocket] error: ${error.message}`);
    console.log(error)

    socket.close();
  }
}

function showConnectionStatus(status) {
  if (!status) {
    return;
  }

  // create the element if it doesn't exist
  if (!document.querySelector('#connection-status')) {
    const div = document.createElement('div');
    div.id = 'connection-status';
    div.style = 'position: fixed; bottom: 0; right: 0; background-color: red; color: white; padding: 5px; font-family: monospace; font-size: 15px; z-index: 9999';
    div.innerText = status;
    document.body.appendChild(div);
  }
}

function deleteConnectionStatus() {
  const element = document.querySelector('#connection-status');
  if (element) {
    element.remove();
  }
}

async function connectionStatusResolver() {
  while (true) {
    showConnectionStatus(wSStatus === 0 ? null : wSStatus === 1 ? 'disconnected' : 'connecting...');
    await new Promise((resolve) => setTimeout(resolve, 100));
    deleteConnectionStatus();
  }
}

async function keepConnectionAlive() {
  while (true) {
    if (wSStatus !== 0) {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      continue;
    }

    // send a ping every 30 seconds to keep the connection alive
    console.log('[websocket] sending ping');
    socket.send('ping');
    await new Promise((resolve) => setTimeout(resolve, 30000));
  }
}

connectionStatusResolver();
connect();
keepConnectionAlive();
