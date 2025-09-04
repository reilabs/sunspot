use std::io::Write;

use tracing::trace;

fn generate_conditional_mov_test(path: &str) {
    let file_name = format!("{path}/conditional_mov.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let conditional_mov = brillig::Opcode::<acir::FieldElement>::ConditionalMov {
        destination: brillig::MemoryAddress::Direct(1234),
        source_a: brillig::MemoryAddress::Direct(5678),
        source_b: brillig::MemoryAddress::Relative(91011),
        condition: brillig::MemoryAddress::Direct(121314),
    }; // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&conditional_mov, config).expect("Failed to encode data");
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
    let directory = format!("{directory}/conditional_mov/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for binary field operations
    generate_conditional_mov_test(&directory);

    trace!("Generated conditional_mov opcode tests in {}", directory);
}
