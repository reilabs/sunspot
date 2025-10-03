use std::hash::poseidon2_permutation;

fn main(input: [Field; 4], hash: pub [Field; 4]) {
    assert(hash == poseidon2_permutation(input, 4))
}

