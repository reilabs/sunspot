fn main(x: u64, y: pub u64, z: u64) {
    assert(x & y == z);
}

#[test]
fn test_main() {
    main(1, 2, 0);

    // Uncomment to make test fail
    // main(1, 1);
}
