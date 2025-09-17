fn main(mut array: [Field; 8], x: u32, y: u32, a: Field, b: Field) {
    assert(pick(array, x) == a);
    array[y] = a;
    assert(array[x] + array[y] == b)
}

fn pick(xs: [Field; 8], idx: u32) -> Field {
   xs[idx]
}
