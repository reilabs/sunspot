use std::{aes128::{aes128_encrypt}};
fn main(inputs: [u8; 16], iv: [u8; 16], key: [u8; 16], ciphertext: [u8;  32]) {
    assert_eq(ciphertext, aes128_encrypt::<16>(inputs, iv, key))
}

