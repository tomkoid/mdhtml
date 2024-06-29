use clap::{Parser, Subcommand};

#[derive(Parser, Default, Clone, Debug)]
#[clap(author = "Tomkoid", version, about = "Converts markdown to html")]
pub struct Args {
    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Subcommand, Clone, Debug)]
pub enum Commands {
    Convert(Convert),
}

#[derive(Parser, Default, Clone, Debug)]
pub struct Convert {
    /// Input file 
    #[clap(required = true)]
    pub input: String,

    /// Output file
    #[clap(short, long, default_value = "[INPUT].html")]
    pub output: String,

    /// Path to a CSS file to use in the HTML
    #[clap(short, long)]
    pub style: Option<String>,

    /// Disables things like highlight.js
    #[clap(short, long, default_value_t = false)]
    pub no_external_libs: bool,

    /// Convert without any additional JavaScript
    #[clap(long)]
    pub raw: bool,

    /// Open in web browser after finishing converting or starting server
    #[clap(short = 'O', long, default_value_t = false)]
    pub open: bool,


    /// Start a HTTP server to serve the HTML file and reload the page when changes are detected
    #[clap(long)]
    pub server: bool,

    /// Hostname of the HTTP server
    #[clap(long, default_value = "127.0.0.1")]
    pub hostname: String,

    /// Port of the HTTP server
    #[clap(long, default_value_t = 3000)]
    pub port: i32,

    /// Prints additional debug info 
    #[clap(long)]
    pub debug: bool,
}
