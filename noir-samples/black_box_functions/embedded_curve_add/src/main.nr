use std::embedded_curve_ops::{embedded_curve_add, EmbeddedCurvePoint};
fn main(x: EmbeddedCurvePoint, y: EmbeddedCurvePoint, z: pub EmbeddedCurvePoint) {

    assert(embedded_curve_add(x, y) == z);
}

