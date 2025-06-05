use tracing::trace;

mod call;
mod memory_init;
mod memory_op;
mod opcode_location;

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/opcodes/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    call::generate_tests(&directory);
    opcode_location::generate_tests(&directory);
    memory_op::generate_tests(&directory);
    memory_init::generate_tests(&directory);

    trace!("Opcode tests generated in {}", directory);
}
