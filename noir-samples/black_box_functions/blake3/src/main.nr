use std::hash::blake3;
fn main(x:[u8; 4], y: [u8; 32] ) {
    let hash= blake3(x);
    assert(hash == y)
}

