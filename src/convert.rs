use std::{
    path::PathBuf,
    process::exit,
    sync::{Arc, Mutex},
};

use colored::Colorize;
use pulldown_cmark::Parser;
use spinners::Spinner;
use tokio::sync::broadcast::{Receiver, Sender};

use crate::{
    http::{message::ChanMessage, server::AppState},
    tools::get_filename,
    SPINNER_TYPE,
};

struct SpinnerState {
    pub stop: bool, // stop signal, stop the spinner if true
}

const SPACING: &str = "   ";

fn spinner_handler(spinner_state: Arc<Mutex<SpinnerState>>) {
    // sleep before starting
    let now = std::time::Instant::now();

    loop {
        if spinner_state.lock().unwrap().stop {
            return;
        }

        if now.elapsed().as_millis() > 200 {
            break;
        }

        std::thread::sleep(std::time::Duration::from_millis(10));
    }

    let mut spinner = Spinner::new(SPINNER_TYPE, "Converting input to markdown...".into());

    loop {
        if spinner_state.lock().unwrap().stop {
            spinner.stop_with_message(format!("Done! Took {}ms.", now.elapsed().as_millis()));
            return;
        }

        std::thread::sleep(std::time::Duration::from_millis(100));
    }
}

pub struct ConvertResponse {
    pub convert_output: Option<String>,
}

pub async fn convert(
    args: &super::args::Convert,
    debug: bool,
    //state: Option<Arc<AppState>>,
    //tx: Option<&Sender<ChanMessage>>
) -> ConvertResponse {
    // config
    let raw_arg = args.raw;
    let style_arg = args.style.clone();
    let server_arg = args.server;

    // get output file name
    let true_output_file = get_filename(&args);

    if server_arg {
        //if state.is_none() {
        //    eprintln!("Error: Could not send transforming message through websocket.");
        //    eprintln!("Could not get state for server.");
        //
        //    exit(1);
        //}

        //Messages::send_transforming_async(axum::extract::State(state.unwrap()));
        //println!("transforming: sending update..");

        //tx.expect("server arg supplied but no tx given").send(ChanMessage {
        //    message: "transforming".to_string(),
        //    status: 2
        //}).unwrap();

        //println!("transforming: sent update!");
    }

    let spinner_state = Arc::new(Mutex::new(SpinnerState { stop: false }));

    let moved_spinner_state = Arc::clone(&spinner_state.clone());
    let spinner_handler_thread = std::thread::spawn(move || spinner_handler(moved_spinner_state));

    let html_output = convert_html_to_markdown(&args.input.to_string());

    let html = if let Ok(html) = html_output {
        html
    } else {
        // sp.stop_with_message("Error converting input to markdown!".into());

        let error: String;

        match html_output.as_ref().unwrap_err().kind() {
            std::io::ErrorKind::NotFound => {
                error = "it doesn't exist.".into();
            }
            std::io::ErrorKind::PermissionDenied => {
                error = "bad permissions.".into();
            }
            std::io::ErrorKind::InvalidData => {
                error = "invalid data.".into();
            }
            std::io::ErrorKind::Other => {
                error = "other error.".into();
            }
            std::io::ErrorKind::UnexpectedEof => {
                error = "unexpected EOF.".into();
            }
            std::io::ErrorKind::AddrInUse => {
                error = "address in use.".into();
            }
            _ => {
                error = html_output.unwrap_err().to_string();
            }
        }

        let message = format!(
            "Error: Could not convert input file {:#?} to markdown, {}.",
            args.input, error
        );

        eprintln!("{}", message);

        exit(1);
    };

    spinner_state.lock().unwrap().stop = true;

    spinner_handler_thread.join().unwrap();

    // read the reload script file
    let reload_script = if server_arg {
        include_str!("http/js/reload.js")
    } else {
        ""
    };

    let scripts = if !args.no_external_libs {
        format!("{}", include_str!("scripts/highlightjs.html"))
    } else {
        "".to_string()
    };

    // styling
    let html = if style_arg.is_some() && !raw_arg {
        let style = std::fs::read_to_string(style_arg.clone().unwrap());

        if style.is_err() {
            eprintln!(
                "Error: Could not read style file {:#?}.",
                style_arg.clone().unwrap()
            );

            eprintln!("{}", style.unwrap_err());

            exit(1);
        }

        let style = style.unwrap();

        format!(
            "<html><head><style>\n{}</style><script>\n{}</script>\n{}</head><body>\n{}</body></html>",
            style, reload_script, scripts, html
        )
    } else {
        if raw_arg {
            html
        } else {
            format!(
                "<html><head><script>{}</script>\n{}</head><body>{}</body></html>",
                reload_script, scripts, html
            )
        }
    };

    // write to file
    let output = std::fs::write(true_output_file.clone(), html);

    if output.is_err() {
        eprintln!(
            "Error: Could not write to output file {:#?}.",
            true_output_file
        );

        eprintln!("{}", output.unwrap_err());

        exit(1);
    }

    if debug {
        println!(
            "{} Successfully wrote to output file {}.",
            "==".green().bold(),
            true_output_file
        );

        let output_path = PathBuf::from(true_output_file)
            .canonicalize()
            .unwrap()
            .display()
            .to_string();

        println!(
            "{SPACING}View in browser at {}{}.",
            "file://".blue(),
                output_path
                .blue()
        );

        if style_arg.is_some() {
            println!(
                "{SPACING}Used style file {}",
                style_arg.clone().unwrap().blue()
            );
        }

        return ConvertResponse {
            convert_output: Some(output_path)
        };
    }

    ConvertResponse {
        convert_output: None
    }
}

fn convert_html_to_markdown(input_file: &String) -> Result<String, std::io::Error> {
    let markdown_input = std::fs::read_to_string(&input_file)?;

    let mut options = pulldown_cmark::Options::empty();
    options.insert(pulldown_cmark::Options::ENABLE_TABLES);
    options.insert(pulldown_cmark::Options::ENABLE_FOOTNOTES);
    options.insert(pulldown_cmark::Options::ENABLE_STRIKETHROUGH);
    options.insert(pulldown_cmark::Options::ENABLE_TASKLISTS);
    options.insert(pulldown_cmark::Options::ENABLE_SMART_PUNCTUATION);
    options.insert(pulldown_cmark::Options::ENABLE_HEADING_ATTRIBUTES);

    let parser = Parser::new_ext(&markdown_input, options);

    let mut html_output = String::new();
    pulldown_cmark::html::push_html(&mut html_output, parser);

    Ok(html_output)
}
