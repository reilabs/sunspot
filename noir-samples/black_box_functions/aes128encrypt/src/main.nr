use std::{aes128::{aes128_encrypt}};
fn main(inputs: [u8; 4]) {
    let iv: [u8; 16] = [0; 16]; // Initialisation vector
    let key: [u8; 16] = [0; 16]; // AES key
    let ciphertext = aes128_encrypt::<4>(inputs, iv, key);
    assert_eq(ciphertext, [181, 52, 14, 195, 227, 220, 39, 229, 171, 134, 202, 77, 26, 234, 135, 83])
}

#[test]
fn test_main() {
   main([0, 12, 3, 15]);

    // Uncomment to make test fail
    // main(1, 1);
}
