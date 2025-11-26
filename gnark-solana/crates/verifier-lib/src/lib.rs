#![warn(missing_docs)]
//! The verifier-lib crate provides utilities for verifying Gnark-generated
//! proofs on Solana.
mod commitments;
mod error;
mod hash;
pub mod proof;
pub mod verifier;
#[cfg(test)]
mod verifier_test;
pub mod vk;
pub mod witness;
