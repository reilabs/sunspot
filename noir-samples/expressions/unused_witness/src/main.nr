fn main(x: Field, y: Field) -> pub Field {
    let mut a = x;
    let mut b = y;
    for _ in 0..5 {
        let s = a + b;
        std::as_witness(s);
        let p = s * s;
        std::as_witness(p);
        a = s;
        b = p;
    }
    a + b
}
