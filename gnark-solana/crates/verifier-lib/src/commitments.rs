//! Utilities for verifying Gnark commitments

use ark_bn254::Fr;
use ark_ff::{BigInteger, PrimeField};
use solana_bn254::prelude::{alt_bn128_g1_multiplication_be, alt_bn128_pairing_be};

use crate::{
    error::GnarkError,
    hash::{hash_to_field, WrappedHashToField},
};

/// Verifies a batched Pedersen proof of knowledge
///
/// The verifier checks the pairing equation
///
/// $$
/// \left(
///   \prod_{i=0}^{n-1} e\left(C_i \cdot \mathrm{challenge}^i, H2_i\right)
/// \right)
/// \cdot
/// e(\mathrm{pok}, G2_{\mathrm{base}})
/// = 1
/// $$
///
/// where
/// - $C_i$ are the Pedersen commitments
/// - $H2_i$ are the corresponding G2 parameters from the verifying key
/// - $\mathrm{challenge}$ is the Fiat–Shamir challenge scalar
/// - $\mathrm{pok}$ is the proof of knowledge element
/// - $G2_{\mathrm{base}}$ is the shared G2 base used for all openings
///
/// The batched proof is accepted only if the full product of pairings
/// equals the identity element of the BN254 target group
pub(crate) fn batch_verify_pedersen(
    vk: &[[u8; 256]],
    commitments: &[[u8; 64]],
    pok: &[u8; 64],
    challenge: Fr,
) -> Result<(), GnarkError> {
    // Ensure parameter sizes and shared G2 base (all vk[i][0..128] equal).
    if commitments.len() != vk.len() {
        return Err(GnarkError::PedersenVerificationError(
            "commitments lengths mismatch".to_string(),
        ));
    }

    for i in 0..vk.len() {
        if i != 0 && vk[i][0..128] != vk[0][0..128] {
            return Err(GnarkError::PedersenVerificationError(
                "parameter mismatch: G2 element".to_string(),
            ));
        }
    }

    // Prepare G1 and G2 pairing inputs:
    //   G1[i] = C_i · challenge^i
    //   G2[i] = H2_i                  (from vk)
    // Final slot:
    //   G1[n] = pok,  G2[n] = G2_base
    let mut pairing_g1: Vec<[u8; 64]> = vec![[0u8; 64]; vk.len() + 1];
    let mut pairing_g2: Vec<[u8; 128]> = vec![[0u8; 128]; vk.len() + 1];

    pairing_g1[0] = commitments[0];
    let mut r = challenge;
    for i in 0..vk.len() {
        let mut arr = [0u8; 128];
        arr.copy_from_slice(&vk[i][128..256]);
        pairing_g2[i] = arr;

        if i != 0 {
            // Compute C_i · challenge^i in G1
            let mut operands = [0u8; 96];
            operands[..64].copy_from_slice(&commitments[i][..]);
            operands[64..96].copy_from_slice(&r.into_bigint().to_bytes_be());
            pairing_g1[i].copy_from_slice(&alt_bn128_g1_multiplication_be(&operands)?);

            if i + 1 != vk.len() {
                r *= challenge;
            }
        }

        pairing_g1[vk.len()] = *pok;
        pairing_g2[vk.len()].copy_from_slice(&vk[0][0..128]);
    }

    let mut pairing_input: Vec<u8> = Vec::with_capacity((64 + 128) * (vk.len() + 1));
    for i in 0..pairing_g1.len() {
        pairing_input.extend_from_slice(&pairing_g1[i]);
        pairing_input.extend_from_slice(&pairing_g2[i]);
    }

    let pairing_res = alt_bn128_pairing_be(pairing_input.as_slice())
        .map_err(|_| GnarkError::ProofVerificationFailed)?;

    // Product of pairings must be the identity
    if pairing_res[31] != 1 {
        return Err(GnarkError::PedersenVerificationError(
            "Pedersen pairing check falied".to_string(),
        ));
    }
    Ok(())
}

/// Computes the Fiat–Shamir challenge scalar used for verifying
/// Gnark commitments.
pub(crate) fn get_challenge<const NR_INPUTS: usize>(
    vk_public_and_commitment_committed: &[&[u64]],
    proof_commitments: &[[u8; 64]],
    public_witness: &mut Vec<[u8; 32]>,
) -> Result<ark_bn254::Fr, GnarkError> {
    let commitments_serialized = solve_commitment_wire::<NR_INPUTS>(
        vk_public_and_commitment_committed,
        proof_commitments,
        public_witness,
    );
    let field_elements = hash_to_field(&commitments_serialized, "G16-BSB22".as_bytes(), 1)?;
    Ok(field_elements[0])
}

/// Computes the derived commitment wire values used in the Gnark commitment
/// scheme.
///
/// For each commitment `i`, the function:
/// 1. Serializes the G1 commitment point `proof_commitments[i]`
/// 2. Appends all public-witness field elements referenced by
///    `vk_public_and_commitment_committed[i]`
/// 3. Hashes the concatenated bytes to a field element using
///    `HashToField("bsb22-commitment")`
/// 4. Appends the resulting field element to `public_witness`
/// 5. Stores the 32-byte hash output in the returned `commitments_serialized`
///    buffer
///
///
/// Mathematically, for each commitment index $i$,
///
/// $$
/// h_i = \mathrm{HashToField}\left(
///   C_i \| \{ x_j \mid j \in \mathrm{indices}(i) \}
/// \right)
/// $$
///
/// where:
/// - $C_i$ is the uncompressed G1 commitment,
/// - $x_j$ are the referenced public-witness field elements,
/// - and $h_i$ is the derived field element appended to the witness.
///
/// Returns a vector containing all serialized $h_i$ values, each encoded in
/// 32 bytes and concatenated in order.
fn solve_commitment_wire<const NR_INPUTS: usize>(
    vk_public_and_commitment_committed: &[&[u64]],
    proof_commitments: &[[u8; 64]],
    public_witness: &mut Vec<[u8; 32]>,
) -> Vec<u8> {
    const FR_BYTES: usize = 32;
    const SIZE_OF_G1_UNCOMPRESSED: usize = 64;

    let mut hash_to_field_fn = WrappedHashToField::new("bsb22-commitment".as_bytes());

    // Compute the maximum number of public committed values in any commitment
    let max_nb_public_committed = vk_public_and_commitment_committed
        .iter()
        .map(|s| s.len())
        .max()
        .unwrap_or(0);

    // Allocate serialized buffers
    let mut commitments_serialized = vec![0u8; vk_public_and_commitment_committed.len() * FR_BYTES];
    let mut commitment_prehash_serialized =
        vec![0u8; SIZE_OF_G1_UNCOMPRESSED + max_nb_public_committed * FR_BYTES];

    for (i, commitment_indices) in vk_public_and_commitment_committed.iter().enumerate() {
        // Copy proof.Commitments[i].Marshal()
        let point_bytes = proof_commitments[i];
        commitment_prehash_serialized[..SIZE_OF_G1_UNCOMPRESSED]
            .copy_from_slice(&point_bytes[..SIZE_OF_G1_UNCOMPRESSED]);

        // Append all corresponding public witness values
        let mut offset = SIZE_OF_G1_UNCOMPRESSED;
        for &j in *commitment_indices {
            let witness_bytes = public_witness[(j - 1) as usize];

            commitment_prehash_serialized[offset..offset + FR_BYTES]
                .copy_from_slice(&witness_bytes);
            offset += FR_BYTES;
        }

        // Hash the values to a single field element and add to public witness and commitments slices
        hash_to_field_fn.write(&commitment_prehash_serialized[..offset]);
        let hash_bts = hash_to_field_fn.sum(Vec::new());
        hash_to_field_fn.reset();

        let mut hash_bts_sized = [0u8; 32];
        hash_bts_sized.copy_from_slice(&hash_bts[0..32]);
        public_witness.push(hash_bts_sized);
        commitments_serialized[i * FR_BYTES..(i + 1) * FR_BYTES].copy_from_slice(&hash_bts);
    }

    commitments_serialized
}
