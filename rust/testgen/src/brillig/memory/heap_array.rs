#![allow(clippy::field_reassign_with_default)]
use std::io::Write;

use brillig::{HeapArray, MemoryAddress};
use tracing::trace;

fn generate_test_heap_array_direct_0x1234(path: &str) {
    let file_name = format!("{path}/heap_array_direct_0x1234.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let mut heap_array = HeapArray::default(); // Example heap array}
    heap_array.pointer = MemoryAddress::Direct(0x1234); // Set the pointer to a direct address
    heap_array.size = 1234; // Set the size of the heap array

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(heap_array, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap array with bytes {:?}",
        file_name, data
    );
}

fn generate_test_heap_array_relative_0x1234(path: &str) {
    let file_name = format!("{path}/heap_array_relative_0x1234.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let mut heap_array = HeapArray::default(); // Example heap array
    heap_array.pointer = MemoryAddress::Relative(0x1234); // Set the pointer to a relative address
    heap_array.size = 1234; // Set the size of the heap array

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(heap_array, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap array with bytes {:?}",
        file_name, data
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/heap_array/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_test_heap_array_direct_0x1234(&directory);
    generate_test_heap_array_relative_0x1234(&directory);

    trace!(
        "Generated heap array test files in directory: {}",
        directory
    );
}
