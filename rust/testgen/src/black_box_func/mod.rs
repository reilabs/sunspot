use std::path::Path;

use tracing::trace;

mod aes128encrypt;
mod and;
mod big_int_add;
mod big_int_div;
mod big_int_from_le_bytes;
mod big_int_mul;
mod big_int_sub;
mod big_int_to_le_bytes;
mod black_box_func_call;
mod black_box_func_kind;
mod blake2s;
mod blake3;
mod ecdsa_secp256k1;
mod ecdsa_secp256r1;
mod embedded_curve_add;
mod function_input;
mod keccak_f1600;
mod multi_scalar_mul;
mod poseidon2permutation;
mod range;
mod recursive_aggregation;
mod sha256_compression;
mod xor;

pub fn generate_tests(root: &str) {
    // Check if the directory exists
    let directory_path = format!("{}/black_box_func", root);
    if !Path::new(&directory_path).exists() {
        // Create the directory
        std::fs::create_dir_all(&directory_path).expect("Failed to create directory");
    }

    aes128encrypt::generate_tests(&directory_path);
    and::generate_tests(&directory_path);
    big_int_add::generate_tests(&directory_path);
    big_int_div::generate_tests(&directory_path);
    big_int_from_le_bytes::generate_tests(&directory_path);
    big_int_mul::generate_tests(&directory_path);
    big_int_sub::generate_tests(&directory_path);
    big_int_to_le_bytes::generate_tests(&directory_path);
    black_box_func_call::generate_tests(&directory_path);
    black_box_func_kind::generate_tests(&directory_path);
    blake2s::generate_tests(&directory_path);
    blake3::generate_tests(&directory_path);
    ecdsa_secp256k1::generate_tests(&directory_path);
    ecdsa_secp256r1::generate_tests(&directory_path);
    embedded_curve_add::generate_tests(&directory_path);
    function_input::generate_tests(&directory_path);
    keccak_f1600::generate_tests(&directory_path);
    multi_scalar_mul::generate_tests(&directory_path);
    poseidon2permutation::generate_tests(&directory_path);
    range::generate_tests(&directory_path);
    recursive_aggregation::generate_tests(&directory_path);
    sha256_compression::generate_tests(&directory_path);
    xor::generate_tests(&directory_path);

    trace!("Generated tests in directory: {}", directory_path);
}
