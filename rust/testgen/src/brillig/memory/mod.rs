mod bitsize;
mod heap_array;
mod heap_value_type;
mod heap_vector;
mod memory_address;
mod value_or_array;

pub fn generate_tests(directory: &str) {
    let directory = format!("{}/memory/", directory);
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    bitsize::generate_tests(&directory);
    heap_array::generate_tests(&directory);
    heap_value_type::generate_tests(&directory);
    heap_vector::generate_tests(&directory);
    memory_address::generate_tests(&directory);
    value_or_array::generate_tests(&directory);
}
