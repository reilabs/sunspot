use std::verify_proof_with_type;

global N: u32 = 18;
global M: u32 = 10;
global K: u32 = 1;

fn main(verification_key: [Field;N], proof: [Field;M], public_inputs: [Field;K], key_hash: Field) {
    let proof_type = 1; 
    verify_proof_with_type(verification_key, proof, public_inputs, key_hash, proof_type)
}

