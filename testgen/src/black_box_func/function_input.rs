use std::io::Write;

use acir::{FieldElement, circuit::opcodes::FunctionInput, native_types::Witness};
use tracing::trace;

fn generate_function_input_test_constant(path: &str) {
    let file_name = format!("{path}/function_input_constant.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let function_input = FunctionInput::<FieldElement>::Constant(FieldElement::from(1234u64));

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(function_input, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_function_input_test_witness(path: &str) {
    let file_name = format!("{path}/function_input_witness.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let function_input = FunctionInput::<FieldElement>::Witness(Witness(1234));

    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(function_input, config).expect("Failed to encode data");
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
    let directory_path = format!("{root}/function_input");
    if !std::path::Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }

    // Generate the tests
    generate_function_input_test_constant(&directory_path);
    generate_function_input_test_witness(&directory_path);

    trace!("Generating tests in directory: {}", directory_path);
}
