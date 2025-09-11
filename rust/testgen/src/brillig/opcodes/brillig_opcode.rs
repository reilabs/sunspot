pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/brillig/opcodes/brillig_opcode/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
}
