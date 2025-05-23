use tracing::trace;

mod witness;

pub fn generate_tests(path: &str) {
    trace!("Running testgen...");
    // Add your test generation logic here
    let directory = format!("{}/shared/", path);

    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");
    witness::generate_witness_tests(&directory);
}
