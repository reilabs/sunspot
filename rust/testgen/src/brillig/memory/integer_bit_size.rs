pub fn generate_tests(directory: &str) {
    let directory = format!("{}/integer_bit_size/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");
}
