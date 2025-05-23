use std::path::Path;

pub fn generate_tests(root: &str) {
    // Check if the directory exists
    let directory_path = format!("{}/black_box_func", root);
    if !Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }
}
