use tracing::trace;

mod expression_test_utils;
mod expression_width;

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/expression/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    expression_width::generate_tests(&directory);
    expression_test_utils::generate_tests(&directory);

    trace!("Expression tests generated in {}", directory);
}
