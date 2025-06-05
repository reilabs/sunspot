use tracing::trace;

pub fn generate_tests(root: &str) {
    // Check if the directory exists
    let directory_path = format!("{}/embedded_curve_add", root);
    if !std::path::Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }

    trace!("Generating tests in directory: {}", directory_path);
}
