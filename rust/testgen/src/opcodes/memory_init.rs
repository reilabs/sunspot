use std::io::Write;

use acir::{
    FieldElement,
    circuit::{
        Opcode,
        opcodes::{BlockId, BlockType},
    },
    native_types::Witness,
};
use tracing::trace;

fn generate_memory_init_memory_block_test(path: &str) {
    let file_name = format!("{}/memory_init_memory_block.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let memory_init_opcode = Opcode::<FieldElement>::MemoryInit {
        block_id: BlockId(0),
        block_type: BlockType::Memory,
        init: vec![],
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&memory_init_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_memory_init_calldata_test(path: &str) {
    let file_name = format!("{}/memory_init_calldata.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let memory_init_opcode = Opcode::<FieldElement>::MemoryInit {
        block_id: BlockId(1),
        block_type: BlockType::CallData(1234),
        init: vec![Witness(0), Witness(1), Witness(2)],
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&memory_init_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generat_memory_init_return_data_test(path: &str) {
    let file_name = format!("{}/memory_init_return_data.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let memory_init_opcode = Opcode::<FieldElement>::MemoryInit {
        block_id: BlockId(2),
        block_type: BlockType::ReturnData,
        init: vec![Witness(0), Witness(1), Witness(2)],
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data =
        bincode::serde::encode_to_vec(&memory_init_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{}/memory_init/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_memory_init_memory_block_test(&directory);
    generate_memory_init_calldata_test(&directory);
    generat_memory_init_return_data_test(&directory);

    trace!("Opcode tests generated in {}", directory);
}
