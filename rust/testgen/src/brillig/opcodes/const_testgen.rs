use std::io::Write;

use acir::FieldElement;
use brillig::{BitSize, IntegerBitSize, MemoryAddress, Opcode};
use tracing::trace;

fn generate_const_tests(path: &str) {
    let file_name = format!("{}/const.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let const_opcode = Opcode::<FieldElement>::Const {
        destination: MemoryAddress::Direct(1234),
        bit_size: BitSize::Integer(IntegerBitSize::U32),
        value: FieldElement::from(5678u32), // Example constant value
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&const_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{}/const/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_const_tests(&directory);

    trace!("Const tests generated in {}", directory);
}
