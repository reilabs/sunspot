use std::{io::Write, path};

use acir::circuit::OpcodeLocation;
use tracing::trace;

fn generate_test_opcode_location_acir(path: &str) {
    let file_name = format!("{}/opcode_location_acir.bin", path);

    // Check if the file already exists
    if path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let opcode_location = OpcodeLocation::Acir(1234);

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&opcode_location, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_test_opcode_location_brillig(path: &str) {
    let file_name = format!("{}/opcode_location_brillig.bin", path);

    // Check if the file already exists
    if path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let opcode_location = OpcodeLocation::Brillig {
        acir_index: 5678,
        brillig_index: 1234,
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&opcode_location, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/opcode_location/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_test_opcode_location_acir(&directory);
    generate_test_opcode_location_brillig(&directory);

    trace!("Opcode tests generated in {}", directory);
}
