use std::io::Write;

use acir::FieldElement;
use brillig::{
    BitSize, HeapArray, HeapValueType, HeapVector, IntegerBitSize, MemoryAddress, Opcode,
    ValueOrArray,
};
use tracing::trace;

fn generate_foreign_call_tests_empty(path: &str) {
    let file_name = format!("{path}/foreign_call_empty.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let foreign_call_opcode = Opcode::<FieldElement>::ForeignCall {
        function: "example_function".to_string(),
        destinations: vec![],
        destination_value_types: vec![],
        inputs: vec![],
        input_value_types: vec![],
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&foreign_call_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_foreign_call_tests_with_inputs(path: &str) {
    let file_name = format!("{path}/foreign_call_with_inputs.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let foreign_call_opcode = Opcode::<FieldElement>::ForeignCall {
        function: "example_function".to_string(),
        destinations: vec![
            ValueOrArray::MemoryAddress(MemoryAddress::Direct(1234)),
            ValueOrArray::HeapArray(HeapArray {
                pointer: MemoryAddress::Direct(5678),
                size: 10,
            }),
        ],
        destination_value_types: vec![HeapValueType::Simple(BitSize::Field)],
        inputs: vec![ValueOrArray::HeapVector(HeapVector {
            pointer: MemoryAddress::Direct(91011),
            size: MemoryAddress::Relative(1212),
        })],
        input_value_types: vec![HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U32))],
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&foreign_call_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{directory}/foreign_call/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_foreign_call_tests_empty(&directory);
    generate_foreign_call_tests_with_inputs(&directory);

    trace!("Foreign call tests generated in {}", directory);
}
