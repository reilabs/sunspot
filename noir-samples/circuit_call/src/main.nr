use std::hash::poseidon2_permutation;
fn main(x: [Field; 4], y: pub [Field; 4]) {
    assert(x != y);
    let z = nested_call(x);
    let z2 = nested_call(y);
    assert(z2 != z);
}

#[fold]
fn nested_call(x: [Field; 4]) -> [Field; 4] {
    inner_call(x)
}

#[fold]
fn inner_call(x: [Field; 4]) -> [Field; 4]{
    poseidon2_permutation(x, 4)
}