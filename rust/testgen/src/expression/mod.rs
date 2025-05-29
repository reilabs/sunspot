use tracing::trace;

mod expression;
mod expression_width;

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/expression/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    expression_width::generate_tests(&directory);
    expression::generate_tests(&directory);    

    trace!("Expression tests generated in {}", directory);
}