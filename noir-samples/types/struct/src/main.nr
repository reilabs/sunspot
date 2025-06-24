struct input {
    a: Field,
    b: Field,
}

fn main(x: input, y: pub Field) {
    assert(x.a + x.b == y);
}

#[test]
fn test_main() {
    main(input {a : 1, b : 1}, 2);

    // Uncomment to make test fail
    // main(1, 1);
}
