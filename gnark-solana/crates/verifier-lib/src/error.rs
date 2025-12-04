// Adapted from https://github.com/Lightprotocol/groth16-solana (Apache 2.0 License)
// Modified by Matthew Klein on 2025-11-26.

use thiserror::Error;

#[derive(Debug, Error, PartialEq)]
pub enum GnarkError {
    #[error("Incompatible Verifying Key with number of public inputs")]
    IncompatibleVerifyingKeyWithNrPublicInputs,
    #[error("ProofVerificationFailed")]
    ProofVerificationFailed,
    #[error("PreparingInputsG1AdditionFailed")]
    PreparingInputsG1AdditionFailed,
    #[error("PreparingInputsG1MulFailed")]
    PreparingInputsG1MulFailed,
    #[error("InvalidG1Length")]
    InvalidG1Length,
    #[error("InvalidG2Length")]
    InvalidG2Length,
    #[error("InvalidPublicInputsLength")]
    InvalidPublicInputsLength,
    #[error("DecompressingG1Failed")]
    DecompressingG1Failed,
    #[error("DecompressingG2Failed")]
    DecompressingG2Failed,
    #[error("PublicInputGreaterThanFieldSize")]
    PublicInputGreaterThanFieldSize,
    #[error("Arkworks serialization error: {0}")]
    ArkworksSerializationError(String),
    #[error("Failed to convert proof component to byte array")]
    ProofConversionError,
    #[error("Failed to compute solana bn254 operation")]
    SolanaBN254Error(String),
    #[error("Error computing FS Hashes")]
    HashError(String),
    #[error("Pedersen verification error")]
    PedersenVerificationError(String),
    #[error("Public witness parsing error")]
    PublicWitnessParsingError,
}

impl From<ark_serialize::SerializationError> for GnarkError {
    fn from(e: ark_serialize::SerializationError) -> Self {
        GnarkError::ArkworksSerializationError(e.to_string())
    }
}

impl From<solana_bn254::AltBn128Error> for GnarkError {
    fn from(e: solana_bn254::AltBn128Error) -> Self {
        GnarkError::SolanaBN254Error(e.to_string())
    }
}

impl From<GnarkError> for u32 {
    fn from(error: GnarkError) -> Self {
        match error {
            GnarkError::IncompatibleVerifyingKeyWithNrPublicInputs => 0,
            GnarkError::ProofVerificationFailed => 1,
            GnarkError::PreparingInputsG1AdditionFailed => 2,
            GnarkError::PreparingInputsG1MulFailed => 3,
            GnarkError::InvalidG1Length => 4,
            GnarkError::InvalidG2Length => 5,
            GnarkError::InvalidPublicInputsLength => 6,
            GnarkError::DecompressingG1Failed => 7,
            GnarkError::DecompressingG2Failed => 8,
            GnarkError::PublicInputGreaterThanFieldSize => 9,
            GnarkError::ArkworksSerializationError(_) => 10,
            GnarkError::ProofConversionError => 11,
            GnarkError::SolanaBN254Error(_) => 12,
            GnarkError::HashError(_) => 13,
            GnarkError::PedersenVerificationError(_) => 14,
            GnarkError::PublicWitnessParsingError => 15,
        }
    }
}
