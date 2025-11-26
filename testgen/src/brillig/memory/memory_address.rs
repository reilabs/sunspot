use std::io::Write;

use brillig::MemoryAddress;
use tracing::trace;

fn generate_test_memory_address_direct_zero(path: &str) {
    let file_name = format!("{path}/memory_address_direct_zero.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let memory_address = MemoryAddress::Direct(0); // Example memory address
    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(memory_address, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for memory address with bytes {:?}",
        file_name, data
    );
}

fn generate_test_memory_address_direct_0x1234(path: &str) {
    let file_name = format!("{path}/memory_address_direct_0x1234.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let memory_address = MemoryAddress::Direct(0x1234); // Example memory address
    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(memory_address, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for memory address with bytes {:?}",
        file_name, data
    );
}

fn generate_test_memory_address_relative_zero(path: &str) {
    let file_name = format!("{path}/memory_address_relative_zero.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let memory_address = MemoryAddress::Relative(0); // Example memory address
    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(memory_address, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for memory address with bytes {:?}",
        file_name, data
    );
}

fn generate_test_memory_address_relative_0x1234(path: &str) {
    let file_name = format!("{path}/memory_address_relative_0x1234.bin");
    // check if the file exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let memory_address = MemoryAddress::Relative(0x1234); // Example memory address
    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(memory_address, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} for memory address with bytes {:?}",
        file_name, data
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/memory_address/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_test_memory_address_direct_zero(&directory);
    generate_test_memory_address_direct_0x1234(&directory);
    generate_test_memory_address_relative_zero(&directory);
    generate_test_memory_address_relative_0x1234(&directory);

    trace!(
        "Generated memory address test files in directory: {}",
        directory
    );
}
