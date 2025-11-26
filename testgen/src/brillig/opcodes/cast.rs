use std::io::Write;

use acir::FieldElement;
use brillig::{BitSize, IntegerBitSize, MemoryAddress, Opcode};
use tracing::trace;

fn generate_cast_test(path: &str) {
    let file_name = format!("{path}/cast.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let cast_opcode = Opcode::<FieldElement>::Cast {
        destination: MemoryAddress::Direct(1234),
        source: MemoryAddress::Direct(5678),
        bit_size: BitSize::Integer(IntegerBitSize::U32),
    }; // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&cast_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{directory}/cast/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_cast_test(&directory);

    trace!("Generated cast test in directory: {}", directory);
}
