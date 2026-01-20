mod binary_field_op;
mod binary_int_op;
mod brillig_opcode;
mod call;
mod call_data_copy;
mod cast;
mod conditional_mov;
mod const_testgen;
mod foreign_call;
mod indirect_const;
mod jump;
mod jump_if;
mod load;
mod mov;
mod not;
mod stop;
mod store;
mod trap;

pub fn generate_tests(directory: &str) {
    let directory = format!("{directory}/opcodes/");
    // Create the directory if it doesn't exist
    std::fs::create_dir_all(&directory).expect("Failed to create directory");

    // Generate witness tests for each opcode module
    binary_field_op::generate_tests(&directory);
    binary_int_op::generate_tests(&directory);
    brillig_opcode::generate_tests(&directory);
    call_data_copy::generate_tests(&directory);
    call::generate_tests(&directory);
    cast::generate_tests(&directory);
    conditional_mov::generate_tests(&directory);
    const_testgen::generate_tests(&directory);
    foreign_call::generate_tests(&directory);
    indirect_const::generate_tests(&directory);
    jump_if::generate_tests(&directory);
    jump::generate_tests(&directory);
    load::generate_tests(&directory);
    mov::generate_tests(&directory);
    not::generate_tests(&directory);
    stop::generate_tests(&directory);
    store::generate_tests(&directory);
    trap::generate_tests(&directory);
}
