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
    #[clap(required = true)]
    pub input: String,
    #[clap(short, long, default_value = "[INPUT].html")]
    pub output: String,
    #[clap(short, long)]
    pub style: Option<String>,
    #[clap(short, long, default_value_t = false)]
    pub no_external_libs: bool,
    #[clap(long)]
    pub raw: bool,

    #[clap(long)]
    pub server: bool,
    #[clap(long)]
    pub debug: bool,
}
