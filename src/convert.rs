use std::{
    path::PathBuf,
    process::exit,
    sync::{Arc, Mutex},
};

use pulldown_cmark::Parser;
use spinners::Spinner;

use crate::{
    http::{message::Messages, server::AppState},
    tools::get_filename,
    SPINNER_TYPE,
};

struct SpinnerState {
    pub stop: bool, // stop signal, stop the spinner if true
}

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

pub fn convert(args: &super::args::Convert, debug: bool, state: Option<AppState>) {
    // config
    let raw_arg = args.raw;
    let style_arg = args.style.clone();
    let server_arg = args.server;

    // get output file name
    let true_output_file = get_filename(&args);

    if server_arg {
        if state.is_none() {
            eprintln!("Error: Could not send transforming message through websocket.");
            eprintln!("Could not get state for server.");

            exit(1);
        }

        Messages::send_transforming_async(axum::extract::State(state.unwrap()));
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
            "<html><head><style>\n{}</style><script>\n{}</script></head><body>\n{}</body></html>",
            style, reload_script, html
        )
    } else {
        if raw_arg {
            html
        } else {
            format!(
                "<html><head><script>{}</script></head><body>{}</body></html>",
                reload_script, html
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
        println!("=> Successfully wrote to output file {}.", true_output_file);

        println!(
            "=> View in browser at file://{}.",
            PathBuf::from(true_output_file)
                .canonicalize()
                .unwrap()
                .display()
        );

        if style_arg.is_some() {
            println!("  => Used style file {}", style_arg.clone().unwrap());
        }
    }
}

fn convert_html_to_markdown(input_file: &String) -> Result<String, std::io::Error> {
    let markdown_input = std::fs::read_to_string(&input_file);

    let markdown_input = match markdown_input {
        Ok(_) => markdown_input.unwrap(),
        Err(e) => {
            return Err(e);
        }
    };

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
