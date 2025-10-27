use age::compare_age;
use date::Date;
use utils::constants::{DG1_LENGTH};


fn main(SAMPLE_DG1: [u8; DG1_LENGTH]) {
     let min_age = 35;
    let max_age = 37;
    let timestamp = Date::new(2025, 9, 12).to_timestamp();
    compare_age(SAMPLE_DG1, min_age, max_age, timestamp);
}

