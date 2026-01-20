use std::io::Write;

use acir::{
    FieldElement,
    circuit::opcodes::{BlackBoxFuncCall, FunctionInput},
    native_types::Witness,
};
use tracing::trace;

fn generate_and_test(path: &str) {
    let file_name = format!("{path}/and_test.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");

    // Example data for the 'and' function call
    let and_function_call = BlackBoxFuncCall::AND {
        lhs: FunctionInput::<FieldElement>::Witness(Witness(1234)),
        rhs: FunctionInput::<FieldElement>::Witness(Witness(2345)),
        num_bits: 5678,
        output: Witness(3456),
    };

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();

    let data =
        bincode::serde::encode_to_vec(&and_function_call, config).expect("Failed to encode data");

    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

pub fn generate_tests(root: &str) {
    // Check if the directory exists
    let directory_path = format!("{root}/and");
    if !std::path::Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }

    generate_and_test(&directory_path);

    trace!("Generating tests in directory: {}", directory_path);
}
