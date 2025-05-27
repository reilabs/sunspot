use std::io::Write;

use brillig::{HeapVector, MemoryAddress};
use tracing::trace;

fn generate_test_heap_vector_zero(path: &str) {
    let file_name = format!("{}/heap_vector_zero.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_vector = HeapVector {
        pointer: MemoryAddress::Direct(0),
        size: MemoryAddress::Relative(0),
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&heap_vector, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap vector with bytes {:?}",
        file_name, data
    );
}

fn generate_test_heap_vector_1234(path: &str) {
    let file_name = format!("{}/heap_vector_1234.bin", path);
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let heap_vector = HeapVector {
        pointer: MemoryAddress::Direct(1234),
        size: MemoryAddress::Relative(5678),
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&heap_vector, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for heap vector with bytes {:?}",
        file_name, data
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/heap_vector/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_test_heap_vector_zero(&directory);
    generate_test_heap_vector_1234(&directory);

    trace!("Generated heap vector tests in directory: {}", directory);
}
