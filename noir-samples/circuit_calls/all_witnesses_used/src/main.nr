
fn main(x: Field, y: pub Field, a: Field, b : pub Field) {
    let z = nested_call(x, y);
    let z2 = nested_call(a, b);
    assert(z + 1 == z2, "z2 must be exactly 1 more than z1");
}

#[fold]
fn nested_call(x: Field, y: Field) -> Field {
    inner_call(x + 2, y)
}

#[fold]
fn inner_call(x: Field, y: Field) -> Field {
    assert(x == y);
    x
}