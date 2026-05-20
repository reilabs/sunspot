#![warn(missing_docs)]
//! The verifier-lib crate provides utilities for verifying Gnark-generated
//! proofs on Solana.
mod commitments;
mod error;
mod hash;
mod proof;
mod syscalls;
mod verifier;
#[cfg(test)]
mod verifier_test;
mod vk;
mod witness;

pub use error::GnarkError;
pub use proof::GnarkProof;
pub use verifier::GnarkVerifier;
pub use vk::{generate_key_file, parse_vk, GnarkVerifyingkey};
pub use witness::GnarkWitness;
