fn main(x: pub Field) {
    let value = x as u8;
    assert(value == 255);
}

#[test]
fn test_main() {
    main(1);

    // Uncomment to make test fail
    // main(1, 1);
}
