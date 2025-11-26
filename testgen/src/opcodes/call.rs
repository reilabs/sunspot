use std::{io::Write, vec};

use acir::{
    FieldElement,
    circuit::{Opcode, opcodes::AcirFunctionId},
    native_types::{Expression, Witness},
};
use tracing::trace;

fn generate_call_test_empty(path: &str) {
    let file_name = format!("{path}/call_empty.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let call_opcode = Opcode::<FieldElement>::Call {
        id: AcirFunctionId(0),
        inputs: vec![],
        outputs: vec![],
        predicate: None,
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&call_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_call_test_with_inputs(path: &str) {
    let file_name = format!("{path}/call_with_inputs.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let call_opcode = Opcode::<FieldElement>::Call {
        id: AcirFunctionId(1),
        inputs: vec![Witness(0), Witness(1), Witness(2), Witness(3), Witness(4)],
        outputs: vec![],
        predicate: None,
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&call_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_call_test_with_outputs(path: &str) {
    let file_name = format!("{path}/call_with_outputs.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let call_opcode = Opcode::<FieldElement>::Call {
        id: AcirFunctionId(2),
        inputs: vec![],
        outputs: vec![Witness(0), Witness(1)],
        predicate: None,
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&call_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_call_test_with_predicate(path: &str) {
    let file_name = format!("{path}/call_with_predicate.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let call_opcode = Opcode::<FieldElement>::Call {
        id: AcirFunctionId(3),
        inputs: vec![],
        outputs: vec![],
        predicate: Some(Expression {
            mul_terms: vec![],
            linear_combinations: vec![],
            q_c: FieldElement::from(1u32),
        }),
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&call_opcode, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_call_test_with_inputs_and_outputs(path: &str) {
    let file_name = format!("{path}/call_with_inputs_and_outputs.bin");

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    let call_opcode = Opcode::<FieldElement>::Call {
        id: AcirFunctionId(4),
        inputs: vec![Witness(0), Witness(1)],
        outputs: vec![Witness(2), Witness(3)],
        predicate: Some(Expression {
            mul_terms: vec![],
            linear_combinations: vec![],
            q_c: FieldElement::from(1u32),
        }),
    };

    // Placeholder for actual data
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&call_opcode, config).expect("Failed to encode data");
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
    let directory = format!("{directory}/call/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_call_test_empty(&directory);
    generate_call_test_with_inputs(&directory);
    generate_call_test_with_outputs(&directory);
    generate_call_test_with_predicate(&directory);
    generate_call_test_with_inputs_and_outputs(&directory);

    trace!("Opcode tests generated in {}", directory);
}
