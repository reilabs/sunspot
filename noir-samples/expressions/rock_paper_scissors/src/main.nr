global N: u32 = 3;

fn main(xs: [Field; N], ys: [Field; N], expected: Field) {
    let mut score: Field = 0;

    for i in 0..N {
        score += xs[i] - ys[i];
    }

    assert(score == expected);
 }
