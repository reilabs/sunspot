use brillig::{HeapArray, HeapVector, MemoryAddress, ValueOrArray};
use std::io::Write;
use tracing::trace;

fn generate_test_value_or_array_memory_address(path: &str) {
    let file_name = format!("{path}/value_or_array_memory_address.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let value_or_array_memory = ValueOrArray::MemoryAddress(MemoryAddress::Direct(1234));

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(value_or_array_memory, config)
        .expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for value or array memory with bytes {:?}",
        file_name, data
    );
}

fn generate_test_value_or_array_heap_array(path: &str) {
    let file_name = format!("{path}/value_or_array_heap_array.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let value_or_array_memory = ValueOrArray::HeapArray(HeapArray {
        pointer: MemoryAddress::Direct(1234),
        size: 5678,
    });

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(value_or_array_memory, config)
        .expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for value or array memory with bytes {:?}",
        file_name, data
    );
}

fn generate_test_value_or_array_heap_vector(path: &str) {
    let file_name = format!("{path}/value_or_array_heap_vector.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let value_or_array_memory = ValueOrArray::HeapVector(HeapVector {
        pointer: MemoryAddress::Direct(1234),
        size: MemoryAddress::Relative(5678),
    });

    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(value_or_array_memory, config)
        .expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for value or array memory with bytes {:?}",
        file_name, data
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/value_or_array/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_test_value_or_array_memory_address(&directory);
    generate_test_value_or_array_heap_array(&directory);
    generate_test_value_or_array_heap_vector(&directory);

    trace!(
        "Generated tests for value or array in directory: {}",
        directory
    );
}
