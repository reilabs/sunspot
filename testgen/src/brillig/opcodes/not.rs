use acir::FieldElement;
use std::io::Write;
use tracing::trace;

use brillig::{IntegerBitSize, MemoryAddress, Opcode};

fn generate_not_test(path: &str) {
    let file_name = format!("{path}/not.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let not_opcode = Opcode::<FieldElement>::Not {
        destination: MemoryAddress::Direct(1234),
        source: MemoryAddress::Relative(5678),
        bit_size: IntegerBitSize::U32,
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&not_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{directory}/not/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_not_test(&directory);

    trace!("Not tests generated in {}", directory);
}
