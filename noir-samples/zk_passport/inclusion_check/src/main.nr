use inclusion_check::check_nationality_inclusion;
use utils::{types::{Alpha3CountryCode}, constants::{EMPTY_ALPHA3_COUNTRY_CODE, DG1_LENGTH}};

fn pad_country_list<let N: u32>(country_list: [Alpha3CountryCode; N]) -> [Alpha3CountryCode; 200] {
    let mut padded_country_list: [Alpha3CountryCode; 200] = [EMPTY_ALPHA3_COUNTRY_CODE; 200];
    for i in 0..N {
        padded_country_list[i] = country_list[i];
    }
    padded_country_list
}

fn main(SAMPLE_DG1: [u8; DG1_LENGTH]) {
     let country_list = pad_country_list(["AUS", "NZL", "USA"]);
    check_nationality_inclusion(SAMPLE_DG1, country_list);
}

