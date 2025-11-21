fn main(mut array: [Field; 8], x: u32, y: u32, a: Field, b: Field) { // copy-on-write
    array[y] = a;           

    let read_x = array[x];
    let read_y = array[y];  

    assert(read_x + read_y == b);
}
