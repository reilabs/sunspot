fn main(x: [Field; 4], y: [Field; 4]) {
    assert(x[0] + x[1] + x[2] + x[3] == y[0] + y[1] + y[2] + y[3]);
}

#[test]
fn test_main() {
    main([1, 2, 3, 4], [2, 3, 3, 2]);

    // Uncomment to make test fail
    // main(1, 1);
}
