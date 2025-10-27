use exclusion_check::check_nationality_exclusion;
use utils::constants::{DG1_LENGTH};


fn main(SAMPLE_DG1: [u8; DG1_LENGTH]) {
    let country_list: [u32; 6] = [4608577, 4608582, 4608590, 4608591, 4608592, 4608595];
    check_nationality_exclusion(SAMPLE_DG1, country_list);
}
