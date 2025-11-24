use utils::constants::DG1_LENGTH;
use sanctions::{types::{SanctionsOrderedMerkleTreeProof, SanctionsOrderedMerkleTreeProofs}, SanctionsSparseMerkleTrees};

fn main(SAMPLE_DG1: [u8; DG1_LENGTH], proofs: [SanctionsOrderedMerkleTreeProof; 4]) {
    let trees = SanctionsSparseMerkleTrees::default();
     let smt_proofs = SanctionsOrderedMerkleTreeProofs {
        // We can use the same proofs for simplicity since there's no second given names
        // they will all be the same
        name_proofs: [proofs[0]; 3],
        passport_no_and_nationality_proof: proofs[1],
        name_and_dob_proofs: [proofs[2]; 3],
        name_and_yob_proofs: [proofs[3]; 3],
    };
    trees.check_sanctions_exclusion(SAMPLE_DG1, smt_proofs, true)
   
}

