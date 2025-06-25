global N: u32 = 3;

fn pow(base: Field, exp: u32) -> Field {
    let mut result: Field = 1;
    for _ in 0..exp {
        result *= base;
    }
    result
}

fn main(coefficients: [Field; N], x : Field, y: pub Field) {
    let mut result: Field = 0;

    for i in 0..N {
        result += coefficients[i] * pow(x, i as u32);
    }

    assert(result == y);
}

#[test]
fn test_main() {
    main([1, 2, 3], 1, 6);

    // Uncomment to make test fail
    // main([1, 2, 3], 1);
}
