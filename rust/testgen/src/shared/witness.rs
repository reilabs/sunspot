use std::io::Write;

use acir::native_types::Witness;
use tracing::trace;

fn generate_zero_witness_test(path: &str) {
    // Check if file exists - if it does, recreate a new one
    let file_path = format!("{}/witness_zero.bin", path);
    if std::path::Path::new(&file_path).exists() {
        std::fs::remove_file(&file_path).expect("Failed to remove existing file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_path).expect("Failed to create file");

    // Write the zero witness to the file
    let witness = Witness::new(0);
    let config = bincode::config::standard()
        .with_little_endian()
        .with_fixed_int_encoding();
    let data = bincode::serde::encode_to_vec(witness, config.clone())
        .expect("msg: Failed to serialize witness");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Zero witness test generated at: {} with data {:?} for {:?}",
        file_path, data, witness
    );
}

fn generate_witness_test_0x1234(path: &str) {
    // Check if file exists - if it does, recreate a new one
    let file_path = format!("{}/witness_1234.bin", path);
    if std::path::Path::new(&file_path).exists() {
        std::fs::remove_file(&file_path).expect("Failed to remove existing file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_path).expect("Failed to create file");

    // Write the witness to the file
    let witness = Witness::new(0x1234);
    let config = bincode::config::standard()
        .with_little_endian()
        .with_fixed_int_encoding();
    let data = bincode::serde::encode_to_vec(witness, config.clone())
        .expect("msg: Failed to serialize witness");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Witness test generated at: {} with data {:?} for {:?}",
        file_path, data, witness
    );
}

fn generate_witness_test_0x12345678(path: &str) {
    // Check if file exists - if it does, recreate a new one
    let file_path = format!("{}/witness_12345678.bin", path);
    if std::path::Path::new(&file_path).exists() {
        std::fs::remove_file(&file_path).expect("Failed to remove existing file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_path).expect("Failed to create file");

    // Write the witness to the file
    let witness = Witness::new(0x12345678);
    let config = bincode::config::standard()
        .with_little_endian()
        .with_fixed_int_encoding();
    let data = bincode::serde::encode_to_vec(witness, config.clone())
        .expect("msg: Failed to serialize witness");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Witness test generated at: {} with data {:?} for {:?}",
        file_path, data, witness
    );
}

pub fn generate_witness_tests(path: &str) {
    let directory = format!("{}/witness/", path);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(directory.clone()).expect("Failed to create directory");
    // Generate zero witness test
    generate_zero_witness_test(directory.as_str());
    generate_witness_test_0x1234(directory.as_str());
    generate_witness_test_0x12345678(directory.as_str());
}
