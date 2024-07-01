use std::sync::Arc;

use axum::{extract::State, response::Html};
use tokio::sync::broadcast::{Receiver, Sender};

use super::server::AppState;

#[derive(Clone, PartialEq, Debug)]
pub struct ChanMessage {
    pub message: String,
    pub status: u16,
    //pub id: i32,
}

// impl Message {
//     pub fn new(message: String, status: u16) -> Self {
//         Self {
//             message,
//             status,
//             id: 0,
//         }
//     }
// }

#[derive(Clone)]
pub struct Messages {
    message_id: i32,
}

impl Messages {
    pub fn new() -> Self {
        Self { message_id: 1 }
    }

    pub fn update(state: Arc<AppState>, message: String, status: u16) {
        let messages_object = &mut state.messages_object.lock().unwrap();

        let message = vec![ChanMessage {
            message: message.into(),
            status: status,
            //id: messages_object.message_id,
        }];

        let current_messages = &mut state.messages.lock().unwrap();

        // add message to messages
        current_messages.extend(message);

        // update message id
        messages_object.message_id += 1;
    }

    // pub fn send_hello(State(state): State<Arc<AppState>>) {
    //     Self::update(state, "hello".into(), 0);
    // }
    //
    pub fn send_update(State(state): State<Arc<AppState>>) {
        Self::update(state, "reload".into(), 1);
    }

    pub fn send_transforming_async(State(state): State<Arc<AppState>>) {
        Self::update(state, "transforming".into(), 2);
    }

    pub async fn send_update_async(State(state): State<Arc<AppState>>) {
        Self::update(state, "reload".into(), 1);
    }

    pub async fn messages_html(State(state): State<Arc<AppState>>) -> Html<String> {
        let mut messages = String::new();

        for message in state.messages.lock().unwrap().iter() {
            messages.push_str(&format!(
                "<div class=\"message {}\">{}</div>",
                match message.status {
                    0 => "info",
                    1 => "success",
                    2 => "error",
                    _ => "info",
                },
                message.message
            ));
        }

        Html(messages)
    }

    pub fn messages(State(state): State<Arc<AppState>>) -> Vec<ChanMessage> {
        state.messages.lock().unwrap().to_vec()
    }
    // pub fn messages_async(State(state): State<Arc<AppState>>) -> Vec<Message> {
    //     state.messages.lock().unwrap().to_vec()
    // }
}

impl Default for ChanMessage {
    fn default() -> Self {
        Self {
            message: String::new(),
            status: 0,
            //id: 0,
        }
    }
}
