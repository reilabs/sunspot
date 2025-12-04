// build.rs
use gnark_verifier_solana::vk::generate_key_file;
use std::env;
use std::path::Path;

fn main() {
    println!("cargo:rerun-if-env-changed=VK_PATH");

    // Default VK path relative to crate root
    let default = "default.vk";
    let vk_path = env::var("VK_PATH").unwrap_or_else(|_| default.to_string());

    if Path::new(&vk_path).exists() {
        println!("cargo:rerun-if-changed={}", vk_path);
    }

    if let Err(e) = generate_key_file(&vk_path, "src/generated_vk.rs") {
        panic!("Failed to generate key file '{}': {e}", vk_path);
    }
}
