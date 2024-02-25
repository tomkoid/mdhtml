use axum::{
    extract::{
        ws::{Message, WebSocket, WebSocketUpgrade},
        State,
    },
    response::IntoResponse,
};

use super::message::Messages;
use super::server::AppState;

pub async fn ws_handler(State(state): State<AppState>, ws: WebSocketUpgrade) -> impl IntoResponse {
    ws.on_upgrade(move |socket| handle_socket(state, socket))
}

async fn handle_socket(state: AppState, mut socket: WebSocket) {
    let state = axum::extract::State(state.clone());
    let mut messages = Messages::messages(state.clone());

    // echo back messages
    socket.send(Message::Text("hello".into())).await.unwrap();

    loop {
        // get messages
        let local_messages = Messages::messages(state.clone());

        // check if messages changed
        if messages != local_messages {
            // get the last id from `messages`
            let last_id = messages
                .last()
                .unwrap_or(&super::message::Message::default())
                .id;

            for message in local_messages.iter().filter(|m| m.id > last_id) {
                let message = socket.send(Message::Text(message.message.clone())).await;

                // check for errors
                if message.is_err() {
                    // client disconnected
                    return;
                }
            }
        }

        tokio::time::sleep(std::time::Duration::from_millis(250)).await;

        // update messages
        messages = local_messages;
    }
}
