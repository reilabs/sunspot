use std::hash::keccak::keccakf1600;

fn main(x: Field, y: pub Field) {
    //assert(x != y);
    let value = keccakf1600([x as u64; 25]);
    assert(value[0] != y as u64);
}

#[test]
fn test_main() {
    main(1, 2);

    // Uncomment to make test fail
    // main(1, 1);
}
