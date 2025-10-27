use expiry::check_expiry;
use date::Date;
use utils::constants::{DG1_LENGTH};

fn main(SAMPLE_DG1: [u8; DG1_LENGTH] ) {
   // This one day before the expiry date, so it is still valid
    let current_date = Date::new(2029, 12, 31);
    check_expiry(SAMPLE_DG1, current_date.to_timestamp());
}

