use std::io::Write;

use acir::acir_field;
use tracing::trace;

fn generate_zero_field(path: &str) {
    let file_name = format!("{path}/zero_field.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let zero_field = acir_field::FieldElement::from(0u64);
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(zero_field, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for zero field with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_1234_field(path: &str) {
    let file_name = format!("{path}/field_1234.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let field_value = acir_field::FieldElement::from(1234u64);
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(field_value, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for field value 1234 with bytes {:?}",
        file_name, data
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/acir_field/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_zero_field(&directory);
    generate_1234_field(&directory);
}
