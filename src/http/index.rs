use std::sync::Arc;

use axum::{extract::State, response::Html};

use crate::tools::get_filename;

use super::server::AppState;

pub async fn get_index(State(state): State<Arc<AppState>>) -> Html<String> {
    // get filename of output file
    let filename = get_filename(&state.args);

    // read converted markdown file
    let index = std::fs::read_to_string(filename);

    // check if file exists
    if index.is_err() {
        return format!(
            "Error: Could not read converted markdown file {:#?}.",
            state.args.output.clone()
        )
        .into();
    }

    // return converted markdown file
    Html(index.unwrap())
}
