use std::{aes128::aes128_encrypt, convert::AsPrimitive};

fn main(data: Field, iv: pub Field, key: pub Field) -> pub [u8; 32] {
    aes128_encrypt([data.as_(); 16], [iv.as_(); 16], [key.as_(); 16])
}

#[test]
fn test_main() {
   //main(1, 2);

    // Uncomment to make test fail
    // main(1, 1);
}
