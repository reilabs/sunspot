use std::io::Write;

use acir::{circuit::opcodes::BlackBoxFuncCall, native_types::Witness, FieldElement};
use tracing::trace;

fn generate_big_int_to_le_bytes_test_empty(path: &str) {
    let file_name = format!("{path}/big_int_to_le_bytes_test.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");

    let big_int_to_le_bytes_function_call = BlackBoxFuncCall::<FieldElement>::BigIntToLeBytes {
        input: 1234,
        outputs: vec![],
    };

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();

    let data = bincode::serde::encode_to_vec(&big_int_to_le_bytes_function_call, config)
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

fn generate_big_int_to_le_bytes_test_with_outputs(path: &str) {
    let file_name = format!("{path}/big_int_to_le_bytes_test_with_outputs.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");

    let big_int_to_le_bytes_function_call = BlackBoxFuncCall::<FieldElement>::BigIntToLeBytes {
        input: 1234,
        outputs: vec![
            Witness(1234),
            Witness(5678),
            Witness(91011),
        ],
    };

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();

    let data = bincode::serde::encode_to_vec(&big_int_to_le_bytes_function_call, config)
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
    let directory_path = format!("{root}/big_int_to_le_bytes");
    if !std::path::Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }

    generate_big_int_to_le_bytes_test_empty(&directory_path);
    generate_big_int_to_le_bytes_test_with_outputs(&directory_path);

    trace!("Generating tests in directory: {}", directory_path);
}
