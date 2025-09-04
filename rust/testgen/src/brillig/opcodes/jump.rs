use std::{io::Write, path};

use acir::FieldElement;
use brillig::{Label, Opcode};
use tracing::trace;

fn generate_jump_test(path: &str) {
    let file_name = format!("{path}/jump.bin");

    // Check if the file already exists
    if path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let jump_opcode = Opcode::<FieldElement>::Jump {
        location: Label::from(1234u16),
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&jump_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{directory}/jump/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_jump_test(&directory);

    trace!("Jump tests generated in {}", directory);
}
