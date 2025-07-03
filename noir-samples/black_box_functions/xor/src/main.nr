use dep::std::convert::AsPrimitive;

fn main(x: Field, y: pub Field, z: pub Field) {
    let lhs: u64 = x.as_();
    let rhs: u64 = y.as_();
    let z: u64 = z.as_();
    assert(lhs ^ rhs == z);
}

#[test]
fn test_main() {
    main(7, 2, 5);

    // Uncomment to make test fail
    // main(1, 1);
}
