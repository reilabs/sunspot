fn main(x: u64, y: u64) {
    assert(x == y);
}

#[test]
fn test_main() {
    main(1, 1);

    // Uncomment to make test fail
    // main(1, 1);
}
