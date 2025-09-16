fn main(x: u64, y: pub u64, z: pub u64) {
    assert(x ^ y == z);
}

#[test]
fn test_main() {
    main(7, 2, 5);

    // Uncomment to make test fail
    // main(1, 1);
}
