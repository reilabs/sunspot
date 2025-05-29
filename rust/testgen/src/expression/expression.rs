use std::io::Write;

use acir::{native_types::{Expression, Witness}, FieldElement};
use tracing::trace;

fn generate_expression_test_empty(path: &str) {
    let file_name = format!("{}/expression_empty.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let expression = Expression::<FieldElement>{
        mul_terms: vec![],
        linear_combinations: vec![],
        q_c : FieldElement::from(0u32),
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    // Serialize the expression to bytes
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&expression, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_expression_test_linear_combinations(path: &str) {
    let file_name = format!("{}/expression_linear_combinations.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let expression = Expression::<FieldElement>{
        mul_terms: vec![],
        linear_combinations: vec![
            (FieldElement::from(1u32), Witness(0)),
            (FieldElement::from(2u32), Witness(1234)),
            (FieldElement::from(3u32), Witness(5678)),
        ],
        q_c : FieldElement::from(0u32),
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    // Serialize the expression to bytes
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&expression, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_expression_test_mul_terms(path: &str) {
    let file_name = format!("{}/expression_mul_terms.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let expression = Expression::<FieldElement>{
        mul_terms: vec![
            (FieldElement::from(1u32), Witness(0), Witness(1)),
            (FieldElement::from(2u32), Witness(1234), Witness(5678)),
            (FieldElement::from(3u32), Witness(5678), Witness(1234)),
        ],
        linear_combinations: vec![],
        q_c : FieldElement::from(0u32),
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    // Serialize the expression to bytes
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&expression, config).expect("Failed to encode data");
    file.write_all(data.as_slice())
        .expect("Failed to write data to file");

    trace!(
        "Generated test file: {} with bytes {:?} len {}",
        file_name,
        data,
        data.len()
    );
}

fn generate_expression_test_mul_terms_with_linear_combinations(path: &str) {
    let file_name = format!("{}/expression_mul_terms_with_linear_combinations.bin", path);

    // Check if the file already exists
    if std::path::Path::new(&file_name).exists() {
        std::fs::remove_file(&file_name).expect("Failed to remove file");
    }

    let expression = Expression::<FieldElement>{
        mul_terms: vec![
            (FieldElement::from(1u32), Witness(0), Witness(1)),
            (FieldElement::from(2u32), Witness(1234), Witness(5678)),
            (FieldElement::from(3u32), Witness(5678), Witness(1234)),
        ],
        linear_combinations: vec![
            (FieldElement::from(1u32), Witness(0)),
            (FieldElement::from(2u32), Witness(1234)),
            (FieldElement::from(3u32), Witness(5678)),
        ],
        q_c : FieldElement::from(0u32),
    };

    // Create a new file
    let mut file = std::fs::File::create(&file_name).expect("Failed to create file");
    // Serialize the expression to bytes
    let config = bincode::config::standard()
        .with_fixed_int_encoding()
        .with_little_endian();
    let data = bincode::serde::encode_to_vec(&expression, config).expect("Failed to encode data");
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
    let directory = format!("{}/expression/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    generate_expression_test_empty(&directory);
    generate_expression_test_linear_combinations(&directory);
    generate_expression_test_mul_terms(&directory);
    generate_expression_test_mul_terms_with_linear_combinations(&directory);

    trace!("Expression tests generated in {}", directory);
}