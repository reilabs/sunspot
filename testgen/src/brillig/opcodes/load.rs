use std::io::Write;

use acir::FieldElement;
use brillig::{MemoryAddress, Opcode};
use tracing::trace;

fn generate_load_test(path: &str) {
    let file_name = format!("{path}/load.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let load_opcode = Opcode::<FieldElement>::Load {
        destination: MemoryAddress::Direct(1234),
        source_pointer: MemoryAddress::Relative(5678),
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&load_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/load/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_load_test(&directory);

    trace!("Load tests generated in {}", directory);
}
