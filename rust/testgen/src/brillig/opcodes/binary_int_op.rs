use std::{io::Write, path};

use acir::FieldElement;
use brillig::{BinaryIntOp, IntegerBitSize, MemoryAddress, Opcode};
use tracing::trace;

fn generate_test_binary_int_op_add(path: &str) {
    let file_name = format!("{}/add.bin", path);

    // Check if the file already exists
    if path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let binary_int_op = Opcode::<FieldElement>::BinaryIntOp {
        destination: MemoryAddress::Direct(1234),
        bit_size: IntegerBitSize::U32,
        op: BinaryIntOp::Add,
        lhs: MemoryAddress::Direct(5678),
        rhs: MemoryAddress::Relative(91011),
    }; // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&binary_int_op, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    tracing::trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_test_binary_int_op_less_than_equals(path: &str) {
    let file_name = format!("{}/less_than_equals.bin", path);

    // Check if the file already exists
    if path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let binary_int_op = Opcode::<FieldElement>::BinaryIntOp {
        destination: MemoryAddress::Direct(123456),
        bit_size: IntegerBitSize::U32,
        op: BinaryIntOp::LessThanEquals,
        lhs: MemoryAddress::Direct(5678910),
        rhs: MemoryAddress::Relative(9101112),
    }; // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&binary_int_op, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    tracing::trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/binary_int_op/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_test_binary_int_op_add(&directory);
    generate_test_binary_int_op_less_than_equals(&directory);

    trace!("Generated binary int op tests in directory: {}", directory);
}
