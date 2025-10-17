fn main(x: Field, y: pub Field) {
    let z = nested_call(x, y);
    let z2 = nested_call(x, y);
    assert(z == z2);
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