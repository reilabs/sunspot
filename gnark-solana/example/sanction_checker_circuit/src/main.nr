use utils::{constants::DG1_LENGTH};
use sanctions::{types::{SanctionsOrderedMerkleTreeProof, SanctionsOrderedMerkleTreeProofs}, SanctionsSparseMerkleTrees};
use commitment_utils::calculate_certificate_registry_leaf;
use validity_check::{check_integrity_of_sod_within_sod_cert, PassportValidityContents};
use commitment_utils::{get_country_from_dg1, compute_merkle_root};
use expiry_check::check_expiry;
use rsa_sig_check::verify_signature;
use tbs_check::verify_rsa_pubkey_in_tbs;


// This circuit certifies that the provided details are a "real" passport
// and that the holder is not on the OFAC sanctions list
fn main(dg1: [u8; DG1_LENGTH], proofs: [SanctionsOrderedMerkleTreeProof; 10], 
    certificate_registry_hash_path: [Field; 16], passport_validity_contents: PassportValidityContents,
    certificate_tags: [Field; 3], certificate_registry_index: Field, current_date: pub u64,
    ) {
    // Registry root from zk_passport (https://registry.zkpassport.id/certificates)
    let certificate_registry_root = 0x2f696abafd61692fe9c82281fd461431f5ff1d3ec31c10b2258b3151d89b9c6d;

    // Make sure that the provided csc is in the merkle tree of certificates
    let leaf: Field =
        calculate_certificate_registry_leaf(certificate_tags, 1, get_country_from_dg1(dg1), passport_validity_contents.csc_pubkey);
    let root = compute_merkle_root(
        leaf,
        certificate_registry_index,
        certificate_registry_hash_path,
    );
    assert(root == certificate_registry_root, "Certificate registry root mismatch");

    // Check that passport is valid and not expired
    check_passport_validity(passport_validity_contents);
    check_expiry(dg1, current_date);

    // Make sure that passport is not in sanction tree
    let trees = SanctionsSparseMerkleTrees::default();
    let smt_proofs = SanctionsOrderedMerkleTreeProofs {
        name_proofs: [proofs[0], proofs[1], proofs[2]],
        passport_no_and_nationality_proof: proofs[3],
        name_and_dob_proofs: [proofs[4], proofs[5], proofs[6]],
        name_and_yob_proofs:  [proofs[7], proofs[8], proofs[9]],
    };
    trees.check_sanctions_exclusion(dg1, smt_proofs, true)
   
}

// CONSTANTS
global SIGNED_ATTRIBUTES_MAX_SIZE: u32 = 200;
global DSC_CERT_MAX_SIZE: u32 = 1300;
global DSC_SIGNATURE_SIZE: u32 = 256;
global CSC_PUBKEY_SIZE: u32 = 512;
fn check_passport_validity(passport_validity_contents: PassportValidityContents) {
    // Checks that H(SOD) is located correctly within the Signed Attributes certificate.
    check_integrity_of_sod_within_sod_cert(
        passport_validity_contents.signed_attributes,
        passport_validity_contents.signed_attributes_size as u32,
        passport_validity_contents.econtent,
        passport_validity_contents.econtent_len,
        passport_validity_contents.econtent_hash_offset,
    );

    // --- SOD certificate signature check ---
    // The second thing asserts verify_sign(message=signed_attributes, sign=sod_signature, pubkey=dsc_pubkey)
    // Using PSS
    assert(verify_signature::<DSC_SIGNATURE_SIZE, 1, SIGNED_ATTRIBUTES_MAX_SIZE, 32>(
        passport_validity_contents.dsc_pubkey,
        passport_validity_contents.dsc_signature,
        passport_validity_contents.dsc_barrett_mu,
        passport_validity_contents.dsc_rsa_exponent,
        passport_validity_contents.signed_attributes,
        passport_validity_contents.signed_attributes_size as u32,
    ));

    // --- DSC certificate signature check ---
    // The first thing checks that the appropriate substring of the `tbs_certificate` is equivalent to the `dsc_pubkey`
    verify_rsa_pubkey_in_tbs(
        passport_validity_contents.dsc_pubkey,
        passport_validity_contents.dsc_cert,
        passport_validity_contents.dsc_pubkey_offset_in_dsc_cert,
    );

    // Check the signature over the DSC
    // Uses SHA-256 with a message size of 1500 and a prime of 2048 bits,
    // using PKCS rather than PSS.
    assert(verify_signature::<CSC_PUBKEY_SIZE, 0, DSC_CERT_MAX_SIZE, 32>(
        passport_validity_contents.csc_pubkey, // This is the pubkey to verify against
        passport_validity_contents.dsc_cert_signature, // This is the actual signature
        passport_validity_contents.csc_barrett_mu,
        passport_validity_contents.csc_rsa_exponent,
        passport_validity_contents.dsc_cert, // This is `data_to_sign`, i.e. the message
        passport_validity_contents.dsc_cert_len,
    ));
}

