fn main(x: Field, y: pub Field, a: pub u64, b: u64) {
    let z1 = foo(a, y, b);
    let z2 = x + z1 as Field;
    assert(z2 == 5);
}

#[fold]
fn foo(x: u64, z: Field, y: u64) -> u64 {
    let r = x & y;
    r
}
