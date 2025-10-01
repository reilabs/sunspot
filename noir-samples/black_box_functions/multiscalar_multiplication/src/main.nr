use std::embedded_curve_ops::{EmbeddedCurvePoint, EmbeddedCurveScalar, multi_scalar_mul};

fn main(points: [EmbeddedCurvePoint; 3], scalars: [EmbeddedCurveScalar; 3], output: EmbeddedCurvePoint) {
    assert(multi_scalar_mul(points, scalars) == output)
}
