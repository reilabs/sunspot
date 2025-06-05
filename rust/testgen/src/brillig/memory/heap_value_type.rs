use std::io::Write;

use brillig::{BitSize, HeapValueType, IntegerBitSize};
use tracing::trace;

fn generate_test_heap_value_type_simple(path: &str) {
    let file_name = format!("{}/heap_value_type_simple.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_value_type = HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U16));

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&heap_value_type, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap value type with bytes {:?}",
        file_name, data
    );
}

fn generate_test_heap_value_type_array_empty(path: &str) {
    let file_name = format!("{}/heap_value_type_array_empty.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_value_type = HeapValueType::Array {
        value_types: vec![],
        size: 10,
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&heap_value_type, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap value type with bytes {:?}",
        file_name, data
    );
}

fn generate_test_heap_value_type_array(path: &str) {
    let file_name = format!("{}/heap_value_type_array.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_value_type = HeapValueType::Array {
        value_types: vec![
            HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U16)),
            HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U32)),
            HeapValueType::Array {
                value_types: vec![
                    HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U8)),
                    HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U16)),
                ],
                size: 123,
            },
        ],
        size: 1234,
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&heap_value_type, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap value type with bytes {:?}",
        file_name, data
    );
}

fn generate_test_heap_value_type_vector_empty(path: &str) {
    let file_name = format!("{}/heap_value_type_vector_empty.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_value_type = HeapValueType::Vector {
        value_types: vec![],
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&heap_value_type, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap value type with bytes {:?}",
        file_name, data
    );
}

fn generate_test_heap_value_type_vector(path: &str) {
    let file_name = format!("{}/heap_value_type_vector.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_value_type = HeapValueType::Vector {
        value_types: vec![
            HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U16)),
            HeapValueType::Simple(BitSize::Integer(IntegerBitSize::U32)),
        ],
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&heap_value_type, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap value type with bytes {:?}",
        file_name, data
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/heap_value_type/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_test_heap_value_type_simple(&directory);
    generate_test_heap_value_type_array_empty(&directory);
    generate_test_heap_value_type_array(&directory);
    generate_test_heap_value_type_vector_empty(&directory);
    generate_test_heap_value_type_vector(&directory);

    trace!(
        "Generated heap value type tests in directory: {}",
        directory
    );
}
