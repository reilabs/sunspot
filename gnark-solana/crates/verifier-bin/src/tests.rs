#[cfg(test)]
mod test {
    use litesvm::LiteSVM;
    use solana_sdk::{
        compute_budget::ComputeBudgetInstruction,
        instruction::Instruction,
        message::Message,
        signature::{Keypair, Signer},
        transaction::Transaction,
    };
    use std::{fs, path::Path};

    #[test]
    fn test_solana_verification() {
        // Create a new LiteSVM instance
        let mut svm = LiteSVM::new();

        // Create a keypair for the transaction payer
        let payer = Keypair::new();

        // Airdrop some lamports to the payer
        svm.airdrop(&payer.pubkey(), 1_000_000_000).unwrap();
        // Load our program
        let program_keypair = Keypair::new();
        let program_id = program_keypair.pubkey();
        svm.add_program_from_file(program_id, "../../target/deploy/verifier_bin.so")
            .unwrap();

        let proof_path = Path::new("test_files/xor.proof");
        let proof_bytes = fs::read(proof_path).expect("Failed to read proof file");
        let witness_path = Path::new("test_files/xor.pw");
        let witness_bytes = fs::read(witness_path).expect("Failed to read witness file");

        let mut data = Vec::with_capacity(proof_bytes.len() + witness_bytes.len());

        data.extend_from_slice(&proof_bytes);
        data.extend_from_slice(&witness_bytes);
        let instruction = Instruction {
            program_id,
            accounts: vec![],
            data,
        };

        // Create transaction
        let compute_budget_ix = ComputeBudgetInstruction::set_compute_unit_limit(500_000); // up to 1.4M CUs
        let message = Message::new(&[compute_budget_ix, instruction], Some(&payer.pubkey()));

        let transaction = Transaction::new(&[&payer], message, svm.latest_blockhash());

        // Send transaction and verify it succeeds
        let result = svm.send_transaction(transaction);
        assert!(result.is_ok(), "Transaction should succeed");
        let logs = result.unwrap().logs;
        println!("Logs: {logs:#?}");
    }

    #[test]
    fn test_solana_verification_should_fail() {
        // Create a new LiteSVM instance
        let mut svm = LiteSVM::new();

        // Create a keypair for the transaction payer
        let payer = Keypair::new();

        // Airdrop some lamports to the payer
        svm.airdrop(&payer.pubkey(), 1_000_000_000).unwrap();
        // Load our program
        let program_keypair = Keypair::new();
        let program_id = program_keypair.pubkey();
        svm.add_program_from_file(program_id, "../../target/deploy/verifier_bin.so")
            .unwrap();

        let proof_path = Path::new("test_files/sum_a_b.proof");
        let proof_bytes = fs::read(proof_path).expect("Failed to read proof file");
        let witness_path = Path::new("test_files/sum_a_b.pw");
        let witness_bytes = fs::read(witness_path).expect("Failed to read witness file");

        let mut data = Vec::with_capacity(proof_bytes.len() + witness_bytes.len());

        data.extend_from_slice(&proof_bytes);
        data.extend_from_slice(&witness_bytes);
        let instruction = Instruction {
            program_id,
            accounts: vec![],
            data,
        };

        // Create transaction
        let compute_budget_ix = ComputeBudgetInstruction::set_compute_unit_limit(500_000); // up to 1.4M CUs
        let message = Message::new(&[compute_budget_ix, instruction], Some(&payer.pubkey()));

        let transaction = Transaction::new(&[&payer], message, svm.latest_blockhash());

        // Send transaction and verify it succeeds
        let result = svm.send_transaction(transaction);
        assert!(result.is_err(), "Transaction should fail");
        let logs = result.unwrap_err().meta.logs;
        println!("Logs: {logs:#?}");
    }
}
