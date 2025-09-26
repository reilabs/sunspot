use std::hash::blake2s; 

fn main(x :[u8; 3], y : [u8; 32]) {
    let hash = blake2s(x);
    assert(hash == y)
}