//! This crate provides a solana verifier program based on the gnark-generated
//! key saved in default.vk, or another file specified by the environment variable VK_PATH.
//! To use another file run `VK_PATH="path_to_vk" cargo build`
mod generated_vk;
mod tests;

use gnark_verifier_solana::{proof::GnarkProof, verifier::GnarkVerifier, witness::GnarkWitness};
use solana_program::{
    account_info::AccountInfo,
    entrypoint::{entrypoint, ProgramResult},
    msg,
    program_error::ProgramError,
    pubkey::Pubkey,
};

entrypoint!(process_instruction);

/// The entrypoint for our program
/// Expects instruction data to be the gnark-generated proof and public witness
/// bytes concatenated together.
/// Will return an error if proof/witness can't be verified
pub fn process_instruction(
    _program_id: &Pubkey,
    _accounts: &[AccountInfo],
    instruction_data: &[u8],
) -> ProgramResult {
    // Number of public inputs
    const NR_INPUTS: usize = generated_vk::VK.nr_pubinputs;
    let proof_len = instruction_data.len() - (12 + NR_INPUTS * 32);
    let proof_bytes = &instruction_data[..proof_len];

    let proof = GnarkProof::from_bytes(proof_bytes).map_err(|e| {
        msg!("Gnark error: {:?}", e);
        ProgramError::Custom(u32::from(e))
    })?;

    let public_witness_bytes = &instruction_data[proof_len..];
    let public_witness = GnarkWitness::from_bytes(public_witness_bytes).map_err(|e| {
        msg!("Gnark error: {:?}", e);
        ProgramError::Custom(u32::from(e))
    })?;

    // Construct the verifier
    let mut verifier: GnarkVerifier<NR_INPUTS> = GnarkVerifier::new(&generated_vk::VK);

    // Perform verification
    let result = verifier.verify(proof, public_witness);

    if result.is_ok() {
        msg!("Proof verified successfully!");
        Ok(())
    } else {
        msg!("Proof verification failed!");
        Err(ProgramError::InvalidInstructionData)
    }
}
