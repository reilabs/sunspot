fn main(x: Field, y: Field, a: pub Field, b: pub Field) {
    assert(a * x + b == y);
}

#[test]
fn test_main() {
    main(1, 2, 1, 0);

    // Uncomment to make test fail
    // main(1, 1);
}
