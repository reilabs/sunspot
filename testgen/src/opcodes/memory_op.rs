use std::io::Write;

use acir::{
    FieldElement,
    circuit::{
        Opcode,
        opcodes::{BlockId, MemOp},
    },
    native_types::Expression,
};
use tracing::trace;

fn generate_memory_op_test_without_predicate(path: &str) {
    let file_name = format!("{path}/memory_op_without_predicate.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let memory_op = Opcode::<FieldElement>::MemoryOp {
        block_id: BlockId(0),
        op: MemOp::<FieldElement> {
            operation: Expression::<FieldElement> {
                mul_terms: vec![],
                linear_combinations: vec![],
                q_c: FieldElement::from(1u32),
            },
            index: Expression::<FieldElement> {
                mul_terms: vec![],
                linear_combinations: vec![],
                q_c: FieldElement::from(2u32),
            },
            value: Expression::<FieldElement> {
                mul_terms: vec![],
                linear_combinations: vec![],
                q_c: FieldElement::from(3u32),
            },
        },
    };
    // Placeholder for actual data

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&memory_op, config)
        .expect("Failed to encode data into bytes");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_memory_op_test_with_predicate(path: &str) {
    let file_name = format!("{path}/memory_op_with_predicate.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let memory_op = Opcode::<FieldElement>::MemoryOp {
        block_id: BlockId(1),
        op: MemOp::<FieldElement> {
            operation: Expression::<FieldElement> {
                mul_terms: vec![],
                linear_combinations: vec![],
                q_c: FieldElement::from(4u32),
            },
            index: Expression::<FieldElement> {
                mul_terms: vec![],
                linear_combinations: vec![],
                q_c: FieldElement::from(5u32),
            },
            value: Expression::<FieldElement> {
                mul_terms: vec![],
                linear_combinations: vec![],
                q_c: FieldElement::from(6u32),
            },
        },
    };
    // Placeholder for actual data

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&memory_op, config)
        .expect("Failed to encode data into bytes");
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
    let directory = format!("{directory}/memory_op/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate the test for the CALL opcode
    generate_memory_op_test_without_predicate(&directory);
    generate_memory_op_test_with_predicate(&directory);

    trace!("Opcode tests generated in {}", directory);
}
