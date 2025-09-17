fn main(mut array: [Field; 8], x: u32, y: u32, a: Field, b: Field) {
    assert(array[x] == a);
    array[y] = a;
    assert(array[x] + array[y] == b)
}