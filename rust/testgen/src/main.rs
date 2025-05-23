use clap::Parser;

mod black_box_func;
mod brillig;
mod shared;

#[derive(Parser, Debug)]
pub struct Params {
    #[clap(short, long, default_value = "../../binaries/")]
    pub target_dir: String,
}

fn main() {
    let params = Params::parse();
    black_box_func::generate_tests(params.target_dir.as_str());

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::TRACE)
        .with_target(false)
        .with_writer(std::io::stdout)
        .init();

    // check if the directory exists
    if !std::path::Path::new(&params.target_dir).exists() {
        std::fs::create_dir_all(&params.target_dir).expect("Failed to create directory");
    }

    shared::generate_tests(params.target_dir.as_str());
    brillig::generate_tests(params.target_dir.as_str());
    println!("Hello, world!");
}
