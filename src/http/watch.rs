use std::path::Path;

use axum::extract::State;
use std::time::Duration;

use crate::{
    args::Convert,
    convert::convert,
    http::message::{ChanMessage, Messages},
};
use sha2::Digest;

use super::server::AppState;

pub async fn watch(args: &Convert, State(state): State<AppState>) -> anyhow::Result<()> {
    println!("Watching for changes... Press Ctrl+C to exit.");

    if args.raw {
        println!("RAW MODE");
    }

    let mut init_hash = String::from("something very random");

    // some error handling
    loop {
        let file_hash = match file_data_hash(&args.input) {
            Ok(hash) => hash,
            Err(e) => {
                if args.debug {
                    eprintln!("{}", e.to_string());
                }

                // sleep for a while
                std::thread::sleep(Duration::from_millis(20));

                continue;
            }
        };

        if file_hash != init_hash {
            convert(&args, false, Some(state.clone())).await;

            let tx = state.tx.lock().await;

            // send messages
            //Messages::send_update(axum::extract::State(state.clone()));
            println!("reload: sending update..");
            tx.send(ChanMessage {
                message: "reload".to_string(),
                status: 1,
            }).unwrap();
            println!("reload: sent update!");

            drop(tx);

            init_hash = file_hash;
        }
    }
}

fn file_data_hash(filename: &str) -> anyhow::Result<String> {
    // check if file exists
    if !Path::new(filename).exists() {
        return Err(anyhow::anyhow!("Could not find file {:#?}", filename));
    }

    // read file
    let file = std::fs::read_to_string(&filename)?;

    // get file hash in sha256
    let file_hash = {
        let mut hasher = sha2::Sha256::new();
        hasher.update(file);
        let result = hasher.finalize();
        format!("{:x}", result)
    };

    // return file hash
    Ok(file_hash)
}
