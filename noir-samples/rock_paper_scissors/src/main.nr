global N: u32 = 3;

fn main(xs: [Field; N], ys: [Field; N], expected: Field) {
    let mut result = 0;
    for i in 0..N {
        let mut score = 0;
        assert((xs[i] as u8) <= 2);
        assert((ys[i] as u8) <= 2);
        let mut score = 0;
        if ((xs[i] as u8) == (ys[i] as u8)) {
            score = 1; // Draw
        } else if (
            (((xs[i] as u8) == 0) * ((ys[i] as u8) == 2)) +
            (((xs[i] as u8) == 1) * ((ys[i] as u8) == 0)) +
            (((xs[i] as u8) == 2) * ((ys[i] as u8) == 1))
            == true
        ) {
            score = 2; // Player wins
        } else {
            score = 0; // Player loses
        }
        result += score;
    }
    assert(result == expected);
 }
