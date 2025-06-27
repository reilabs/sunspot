use std::hash::sha256;

fn main(data1: Field, data2: pub Field) {
    // Convert the Field to bytes (assuming 16 bytes for AES128, adjust as needed)
    let data_bytes1: [u8; 16] = data1.to_le_bytes();
    let hash1 = sha256(data_bytes1);
    let data_bytes2: [u8; 16] = data2.to_le_bytes();
    let hash2 = sha256(data_bytes2);
    assert_eq(hash1, hash2);
}

#[test]
fn test_main() {
   //main(1, 2);

    // Uncomment to make test fail
    // main(1, 1);
}
