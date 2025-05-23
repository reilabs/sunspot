pub fn generate_tests(directory: &str) {
    let directory = format!("{}/heap_array/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");
}
