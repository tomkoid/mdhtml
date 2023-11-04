function connect() {
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

    // send a ping every 30 seconds to keep the connection alive
    setInterval(() => {
      console.log('[websocket] sending ping');
      socket.send('ping');
    }, 30000);
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
    console.log('[websocket] connection closed');

    alert("Connection to server lost. Please refresh the page to reconnect.")
  }

  socket.onerror = (error) => {
    console.log(`[websocket] error: ${error.message}`);
    console.log(error)

    socket.close();
  }
}

connect();
