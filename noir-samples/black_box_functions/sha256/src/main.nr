fn main(x: Field, y: pub Field) {
    //let x : [u8] = x.to_be_bytes();
    //let hash = sha256::sha256_var(x, "sha256");
    assert(x != y);
}

#[test]
fn test_main() {
    main(1, 2);

    // Uncomment to make test fail
    // main(1, 1);
}
