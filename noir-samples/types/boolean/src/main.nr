fn main(x: bool, y: pub bool) {
    assert(x != y);
}

#[test]
fn test_main() {
    main(true, false);
}
