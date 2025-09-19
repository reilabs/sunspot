use sha256::{sha256_var};
fn main(input: [u8;1], result: [u8; 32]) {
    assert( sha256_var(input, 1) == result)
}

