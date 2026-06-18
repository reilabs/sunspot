use std::io::Write;

use crate::encode;
use acir::{
    FieldElement,
    circuit::{
        Opcode,
        opcodes::{BlockId, MemOp, MemOpKind},
    },
    native_types::Witness,
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
        op: MemOp {
            operation: MemOpKind::Read,
            index: Witness(2),
            value: Witness(3),
        },
    };
    // Placeholder for actual data

    let data = encode(&memory_op);
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
        op: MemOp {
            operation: MemOpKind::Write,
            index: Witness(5),
            value: Witness(6),
        },
    };
    // Placeholder for actual data

    let data = encode(&memory_op);
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
