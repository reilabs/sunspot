use std::hash::blake2s; 

fn main(x :[u8; 3], y : [u8; 32], a:[u8; 128], b: [u8; 32] ) {
    assert(blake2s(x) == y);
    assert(blake2s(a) == b)
}