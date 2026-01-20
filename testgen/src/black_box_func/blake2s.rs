use std::io::Write;

use acir::{
    FieldElement,
    circuit::opcodes::{BlackBoxFuncCall, FunctionInput},
    native_types::Witness,
};
use tracing::trace;

fn generate_blake2s_test_empty(path: &str) {
    let file_name = format!("{path}/blake2s_test_empty.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");

    let blake2s_function_call = BlackBoxFuncCall::<FieldElement>::Blake2s {
        inputs: vec![],
        outputs: Box::new([Witness(0); 32]),
    };

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();

    let data = bincode::serde::encode_to_vec(&blake2s_function_call, config)
        .expect("Failed to encode data");

    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_blake2s_test_with_inputs(path: &str) {
    let file_name = format!("{path}/blake2s_test_with_inputs.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");

    let blake2s_function_call = BlackBoxFuncCall::<FieldElement>::Blake2s {
        inputs: vec![
            FunctionInput::Witness(Witness(1234)),
            FunctionInput::Witness(Witness(5678)),
        ],
        outputs: Box::new([Witness(1234); 32]),
    };

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();

    let data = bincode::serde::encode_to_vec(&blake2s_function_call, config)
        .expect("Failed to encode data");

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
    let directory_path = format!("{root}/blake2s");
    if !std::path::Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }

    generate_blake2s_test_empty(&directory_path);
    generate_blake2s_test_with_inputs(&directory_path);

    trace!("Generating tests in directory: {}", directory_path);
}
