fn main(x: Field, a: pub Field, b: pub Field, c: pub Field) {
    assert(a * x * x + b * x + c == 0);
}

#[test]
fn test_main() {
    main(-1, 1, 2, 1);
}
