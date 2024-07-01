use std::{sync::Arc, time::Duration};

use axum::{
    extract::{
        ws::{Message, WebSocket, WebSocketUpgrade},
        State,
    },
    response::IntoResponse,
};
use tokio::sync::broadcast::Receiver;
use tokio_stream::{wrappers::BroadcastStream, StreamExt};

use crate::http::message::ChanMessage;

use super::server::AppState;

pub async fn ws_handler(State(state): State<Arc<AppState>>, ws: WebSocketUpgrade) -> impl IntoResponse {
    ws.on_upgrade(move |socket| handle_socket(state, socket))
}

async fn handle_socket(state: Arc<AppState>, mut socket: WebSocket) {
    println!("ws: hello client!!");
    let state = axum::extract::State(state.clone());

    // echo back messages
    socket.send(Message::Text("hello".into())).await.unwrap();

    let mut rx: Receiver<ChanMessage> = state.tx.subscribe();

    loop {
        if rx.is_empty() {
            std::thread::sleep(Duration::from_millis(20));
            continue
        }
        println!("waiting for msg..");
        let msg = rx.try_recv();

        let msg = if let Ok(msg) = msg {
            msg
        } else {
            eprintln!("receiver failed: {}", msg.unwrap_err());
            continue;
        };

        println!("msg arrived!");
        println!("chan msg: {}", msg.message);
        let socket_send = socket.send(Message::Text(msg.message)).await;

        if !socket_send.is_ok() {
            eprintln!("error: {}", socket_send.unwrap_err().to_string());
            socket.close().await.unwrap();
            break;
        }

        println!("waiting for msg..");
    }

    //loop {
    // get messages
    //let local_messages = Messages::messages(state.clone());

    // check if messages changed
    //if messages != local_messages {
    //    // get the last id from `messages`
    //    let last_id = messages
    //        .last()
    //        .unwrap_or(&super::message::Message::default())
    //        .id;
    //
    //    for message in local_messages.iter().filter(|m| m.id > last_id) {
    //        let message = socket.send(Message::Text(message.message.clone())).await;
    //
    //        // check for errors
    //        if message.is_err() {
    //            // client disconnected
    //            return;
    //        }
    //    }
    //}

    //tokio::time::sleep(std::time::Duration::from_millis(250)).await;

    // update messages
    //messages = local_messages;
    //}
}
