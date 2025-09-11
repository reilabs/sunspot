mod acir_field;
mod black_box_func;
mod brillig;
mod expression;
mod opcodes;
mod shared;

fn main() {
    let target_dir = "../../binaries/";

    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::TRACE)
        .with_target(false)
        .with_writer(std::io::stdout)
        .init();

    // check if the directory exists
    if !std::path::Path::new(target_dir).exists() {
        std::fs::create_dir_all(target_dir).expect("Failed to create directory");
    }

    black_box_func::generate_tests(target_dir);
    shared::generate_tests(target_dir);
    brillig::generate_tests(target_dir);
    acir_field::generate_tests(target_dir);
    expression::generate_tests(target_dir);
    opcodes::generate_tests(target_dir);

    println!("Hello, world!");
}
