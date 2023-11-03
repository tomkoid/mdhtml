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
        const newScript = document.createElement('script');
        newScript.src = script.src;
        newScript.async = false;
        document.body.appendChild(newScript);
      });
    }
  }

  socket.onclose = () => {
    console.log('[websocket] disconnected');
    console.log("reconnecting...")

    socket.close();

    setTimeout(() => {
      connect();
    }, 1000);
  }

  socket.onerror = (error) => {
    console.log(`[websocket] error: ${error}`);
    console.log(error)

    // this will go to the onclose function
    socket.close();
  }
}

connect();
