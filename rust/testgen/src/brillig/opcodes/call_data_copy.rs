use std::io::Write;

use acir::FieldElement;
use brillig::{MemoryAddress, Opcode};
use tracing::trace;

fn generate_call_data_copy_test(path: &str) {
    let file_name = format!("{}/call_data_copy.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let call_data_copy = Opcode::<FieldElement>::CalldataCopy {
        destination_address: MemoryAddress::Direct(1234),
        size_address: MemoryAddress::Direct(5678),
        offset_address: MemoryAddress::Relative(91011),
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&call_data_copy, config).expect("Failed to encode data");
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
    let directory = format!("{}/call_data_copy/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_call_data_copy_test(&directory);

    trace!("Generated call_data_copy test in directory: {}", directory);
}
