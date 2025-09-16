fn main(x: [Field; 8], y: u32, z: Field) {
    assert(pick(x, y) ==z)
}

fn pick(xs: [Field; 8], idx: u32) -> Field {
   xs[idx]
}


#[test]
fn test_main() {
    main([1,2,3,4,5,6,7,8],  7, 8)
}
