use thiserror::Error;

#[derive(Debug, Error, PartialEq)]
pub enum Groth16Error {
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

impl From<ark_serialize::SerializationError> for Groth16Error {
    fn from(e: ark_serialize::SerializationError) -> Self {
        Groth16Error::ArkworksSerializationError(e.to_string())
    }
}

impl From<solana_bn254::AltBn128Error> for Groth16Error {
    fn from(e: solana_bn254::AltBn128Error) -> Self {
        Groth16Error::SolanaBN254Error(e.to_string())
    }
}

impl From<Groth16Error> for u32 {
    fn from(error: Groth16Error) -> Self {
        match error {
            Groth16Error::IncompatibleVerifyingKeyWithNrPublicInputs => 0,
            Groth16Error::ProofVerificationFailed => 1,
            Groth16Error::PreparingInputsG1AdditionFailed => 2,
            Groth16Error::PreparingInputsG1MulFailed => 3,
            Groth16Error::InvalidG1Length => 4,
            Groth16Error::InvalidG2Length => 5,
            Groth16Error::InvalidPublicInputsLength => 6,
            Groth16Error::DecompressingG1Failed => 7,
            Groth16Error::DecompressingG2Failed => 8,
            Groth16Error::PublicInputGreaterThanFieldSize => 9,
            Groth16Error::ArkworksSerializationError(_) => 10,
            Groth16Error::ProofConversionError => 11,
            Groth16Error::SolanaBN254Error(_) => 12,
            Groth16Error::HashError(_) => 13,
            Groth16Error::PedersenVerificationError(_) => 14,
            Groth16Error::PublicWitnessParsingError => 15,
        }
    }
}
