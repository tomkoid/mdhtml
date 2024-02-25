use std::sync::{Arc, Mutex};

use axum::{routing::get, Router};

use crate::args::Convert;

use super::{
    message::{Message, Messages},
    watch::watch,
};

#[derive(Clone)]
pub struct AppState {
    pub args: Convert,
    pub messages: Arc<Mutex<Vec<Message>>>,
    pub messages_object: Arc<Mutex<Messages>>,
}

pub async fn start_server(args: &Convert) {
    let messages = Messages::new();

    let state = AppState {
        args: args.clone(),
        messages: Arc::new(Mutex::new(Vec::new())),
        messages_object: Arc::new(Mutex::new(messages.clone())),
    };

    // build our application with a single route
    let app = Router::new()
        .route("/", get(super::index::get_index))
        .route("/update", get(super::message::Messages::send_update_async))
        .route("/messages", get(super::message::Messages::messages_html))
        .route("/ws", get(super::websocket::ws_handler))
        .with_state(state.clone());

    let args_temp = args.clone();

    tokio::spawn(async move {
        watch(&args_temp, axum::extract::State(state)).unwrap();
    });

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
