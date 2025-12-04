mod vk;
use gnark_verifier_solana::{proof::GnarkProof, verifier::GnarkVerifier, witness::GnarkWitness};
use solana_program::{
    account_info::{next_account_info, AccountInfo},
    entrypoint::{entrypoint, ProgramResult},
    msg,
    program::invoke_signed,
    program_error::ProgramError,
    program_pack::Pack,
    pubkey::Pubkey,
};
use spl_token::instruction::mint_to;

// Custom condition error
#[derive(Debug)]
pub enum MintingError {
    NotEligible,
}

impl From<MintingError> for ProgramError {
    fn from(e: MintingError) -> Self {
        match e {
            MintingError::NotEligible => ProgramError::Custom(0),
        }
    }
}

entrypoint!(process_instruction);

pub fn process_instruction(
    program_id: &Pubkey,
    accounts: &[AccountInfo],
    instruction_data: &[u8],
) -> ProgramResult {
    let accounts_iter = &mut accounts.iter();

    let mint_account = next_account_info(accounts_iter)?;
    let destination = next_account_info(accounts_iter)?;
    let token_program = next_account_info(accounts_iter)?;
    let mint_authority = next_account_info(accounts_iter)?;

    // Construct the proof and witness and verify
    const NR_INPUTS: usize = vk::VK.nr_pubinputs;
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

    let mut verifier: GnarkVerifier<NR_INPUTS> = GnarkVerifier::new(&vk::VK);
    let result = verifier.verify(proof, public_witness);

    if result.is_err() {
        msg!("User not eligible, unsanctioned status not verified");
        return Err(MintingError::NotEligible.into());
    }

    // Check destination token balance
    let destination_data = spl_token::state::Account::unpack(&destination.try_borrow_data()?)?;
    if destination_data.amount > 0 {
        msg!("Destination account already holds tokens, minting denied");
        return Err(MintingError::NotEligible.into());
    }

    // If eligible and balance is 0, mint tokens
    let amount_to_mint = 1u64;
    msg!("User eligible, minting tokens...");

    let (mint_authority_pda, bump) =
        Pubkey::find_program_address(&[b"ofac_check_mint_auth"], program_id);
    msg!("{}", mint_authority_pda);
    let signer_seeds: &[&[u8]] = &[b"ofac_check_mint_auth", &[bump]];

    invoke_signed(
        &mint_to(
            token_program.key,
            mint_account.key,
            destination.key,
            &mint_authority_pda,
            &[],
            amount_to_mint,
        )?,
        &[
            mint_account.clone(),
            destination.clone(),
            token_program.clone(),
            mint_authority.clone(),
        ],
        &[signer_seeds],
    )?;

    msg!("Minted {} token to {}", amount_to_mint, destination.key);
    Ok(())
}
