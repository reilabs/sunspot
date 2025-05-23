
use std::io::Write;

use brillig::{BitSize, IntegerBitSize};
use tracing::trace;

fn generate_test_bitsize_field(path: &str) {
    let file_name = format!("{}/bitsize_field.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let bitsize = BitSize::Field;
    let config = bincode::config::standard().with_fixed_int_encoding().with_little_endian();
    let data = bincode::serde::encode_to_vec(&bitsize, config).expect("Failed to encode data");
    file.write_all(data.as_slice()).expect("Failed to write data to file");

    trace!("Generated test file: {} for bitsize {:?} with bytes {:?}", file_name, bitsize, data);
}

fn generate_test_bitsize_integer_u1(path: &str) {
    let file_name = format!("{}/bitsize_integer_zero.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let bitsize = BitSize::Integer(IntegerBitSize::U1);
    let config = bincode::config::standard().with_fixed_int_encoding().with_little_endian();
    let data = bincode::serde::encode_to_vec(&bitsize, config).expect("Failed to encode data");
    file.write_all(data.as_slice()).expect("Failed to write data to file");

    trace!("Generated test file: {} for bitsize {:?} with bytes {:?}", file_name, bitsize, data);
}

fn generate_test_bitsize_integer_u8(path: &str) {
    let file_name = format!("{}/bitsize_integer_0x1234.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let bitsize = BitSize::Integer(IntegerBitSize::U8);
    let config = bincode::config::standard().with_fixed_int_encoding().with_little_endian();
    let data = bincode::serde::encode_to_vec(&bitsize, config).expect("Failed to encode data");
    file.write_all(data.as_slice()).expect("Failed to write data to file");

    trace!("Generated test file: {} for bitsize {:?} with bytes {:?}", file_name, bitsize, data);
}

fn generate_test_bitsize_integer_u128(path: &str) {
    let file_name = format!("{}/bitsize_integer_0x123456.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let bitsize = BitSize::Integer(IntegerBitSize::U128);
    let config = bincode::config::standard().with_fixed_int_encoding().with_little_endian();
    let data = bincode::serde::encode_to_vec(&bitsize, config).expect("Failed to encode data");
    file.write_all(data.as_slice()).expect("Failed to write data to file");

    trace!("Generated test file: {} for bitsize {:?} with bytes {:?}", file_name, bitsize, data);
}


pub fn generate_tests(directory: &str) {
    let directory = format!("{}/bitsize/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");
    generate_test_bitsize_field(&directory);
    generate_test_bitsize_integer_u1(&directory);
    generate_test_bitsize_integer_u8(&directory);
    generate_test_bitsize_integer_u128(&directory);
}
