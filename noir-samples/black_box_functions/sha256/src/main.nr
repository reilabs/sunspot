use std::hash::sha256_compression;
fn main(input: [u32; 16], state: [u32; 8], y: [u32; 8]) {
    let output = sha256_compression(input, state);
    assert(output == y)
  
    
}