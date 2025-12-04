//! Provides utilities for parsing Gnark-generated proofs
use std::io::{self, Read};

use crate::{error::GnarkError, vk::read_vk_ic};

/// The Gnark elliptic curve proof elements.
/// Notation follows Figure 4. in DIZK paper <https://eprint.iacr.org/2018/691.pdf>
pub struct GnarkProof<'a> {
    /// G1 element
    pub ar: [u8; 64],
    /// G2 element
    pub bs: [u8; 128],
    /// G1 element
    pub krs: [u8; 64],
    /// Pedersen commitments a la <https://eprint.iacr.org/2022/1072>
    pub commitments: &'a [[u8; 64]],
    /// Batched proof of knowledge of the above commitments
    pub commitment_pok: [u8; 64],
}

impl GnarkProof<'_> {
    /// Parses the Gnark proof from a reader.
    pub fn parse<R: Read>(mut reader: R) -> io::Result<Self> {
        let mut proof_a = [0u8; 64];
        reader.read_exact(&mut proof_a)?;

        let mut proof_b = [0u8; 128];
        reader.read_exact(&mut proof_b)?;

        let mut proof_c = [0u8; 64];
        reader.read_exact(&mut proof_c)?;

        let commitments_vec = read_vk_ic(&mut reader)?;
        let commitments: &'static [[u8; 64]] = Box::leak(commitments_vec.into_boxed_slice());

        let mut commitment_pok = [0u8; 64];
        reader.read_exact(&mut commitment_pok)?;

        Ok(Self {
            ar: proof_a,
            bs: proof_b,
            krs: proof_c,
            commitments,
            commitment_pok,
        })
    }
    /// Parses the Groth16 proof from bytes.
    /// Should be of length 324 + N_COMMITMENTS * 64
    pub fn from_bytes(bytes: &[u8]) -> Result<Self, GnarkError> {
        if bytes.len() < 256 + 4 + 64 {
            return Err(GnarkError::ProofConversionError);
        }

        // Parse A, B, C
        let mut proof_a = [0u8; 64];
        let mut proof_b = [0u8; 128];
        let mut proof_c = [0u8; 64];

        proof_a.copy_from_slice(&bytes[0..64]);
        proof_b.copy_from_slice(&bytes[64..192]);
        proof_c.copy_from_slice(&bytes[192..256]);

        // Parse commitments
        let mut offset = 256;

        // First 4 bytes = number of commitments (u32)
        let num_commitments = u32::from_be_bytes(
            bytes[offset..offset + 4]
                .try_into()
                .map_err(|_| GnarkError::ProofConversionError)?,
        ) as usize;
        offset += 4;

        let expected_len = offset + num_commitments * 64 + 64;
        if bytes.len() != expected_len {
            return Err(GnarkError::ProofConversionError);
        }

        let mut commitments_vec = Vec::with_capacity(num_commitments);
        for _ in 0..num_commitments {
            let mut c = [0u8; 64];
            c.copy_from_slice(&bytes[offset..offset + 64]);
            offset += 64;
            commitments_vec.push(c);
        }

        let commitments: &'static [[u8; 64]] = Box::leak(commitments_vec.into_boxed_slice());

        // Parse commitment_pok
        let mut commitment_pok = [0u8; 64];
        commitment_pok.copy_from_slice(&bytes[offset..offset + 64]);

        Ok(Self {
            ar: proof_a,
            bs: proof_b,
            krs: proof_c,
            commitments,
            commitment_pok,
        })
    }
}

#[cfg(test)]
mod tests {
    use std::fs::File;

    #[test]
    fn test_parse_proof_no_commitment() {
        // Open the test file
        let file = File::open("src/test_files/sum_a_b.proof").unwrap();

        // Parse the verifying key
        let proof = super::GnarkProof::parse(file);

        assert!(proof.is_ok())
    }

    #[test]
    fn test_parse_proof_with_commitment() {
        // Open the test file
        let file = File::open("src/test_files/keccak_f1600.proof").unwrap();

        // Parse the verifying key
        let proof = super::GnarkProof::parse(file);
        assert!(proof.is_ok())
    }

    #[test]
    fn test_proof_from_bytes() {
        // Open the test file
        let bytes = [
            12, 13, 185, 113, 204, 181, 32, 215, 140, 23, 106, 112, 13, 51, 186, 118, 103, 22, 127,
            118, 186, 157, 23, 93, 183, 160, 80, 50, 164, 212, 124, 95, 8, 233, 167, 48, 107, 176,
            228, 108, 238, 254, 124, 226, 207, 72, 172, 95, 254, 156, 203, 154, 129, 148, 124, 96,
            167, 143, 89, 96, 83, 67, 17, 158, 7, 173, 37, 140, 199, 178, 136, 171, 151, 240, 28,
            255, 142, 41, 164, 107, 200, 200, 201, 165, 11, 14, 222, 97, 84, 130, 114, 155, 142,
            220, 125, 202, 20, 177, 223, 211, 160, 187, 66, 211, 84, 133, 152, 2, 83, 51, 76, 237,
            181, 29, 139, 163, 154, 142, 10, 49, 2, 144, 211, 182, 5, 9, 63, 198, 35, 130, 241, 48,
            6, 128, 2, 246, 104, 59, 129, 23, 201, 8, 155, 251, 254, 176, 41, 53, 85, 228, 121, 74,
            202, 201, 159, 240, 36, 7, 5, 59, 4, 23, 113, 24, 119, 118, 222, 218, 118, 232, 133,
            30, 66, 137, 48, 49, 217, 108, 175, 175, 192, 137, 120, 184, 242, 188, 147, 212, 32,
            124, 74, 125, 45, 79, 156, 106, 190, 83, 104, 34, 5, 25, 165, 7, 173, 73, 132, 139, 22,
            178, 148, 5, 111, 92, 87, 31, 2, 122, 228, 144, 3, 196, 117, 5, 17, 228, 246, 211, 220,
            59, 96, 111, 101, 251, 209, 162, 142, 222, 111, 197, 247, 102, 67, 74, 104, 91, 63,
            142, 121, 13, 92, 191, 232, 195, 94, 30, 0, 0, 0, 1, 46, 71, 5, 84, 144, 44, 118, 168,
            148, 30, 15, 4, 229, 145, 155, 94, 235, 167, 2, 206, 113, 71, 3, 179, 41, 95, 166, 223,
            124, 121, 64, 11, 41, 26, 203, 114, 227, 35, 77, 72, 156, 240, 241, 123, 13, 170, 172,
            253, 31, 47, 25, 224, 32, 88, 218, 231, 89, 228, 56, 102, 10, 182, 48, 151, 1, 174,
            234, 110, 183, 215, 224, 79, 199, 78, 132, 214, 253, 81, 159, 28, 221, 228, 126, 195,
            122, 82, 201, 39, 201, 0, 112, 89, 20, 104, 69, 201, 8, 164, 36, 9, 148, 93, 75, 3,
            152, 49, 65, 136, 84, 150, 201, 95, 12, 91, 205, 170, 25, 133, 207, 15, 251, 136, 239,
            50, 32, 186, 121, 113,
        ];

        let proof = super::GnarkProof::from_bytes(&bytes);
        assert!(proof.is_ok())
    }
}
