//! Provides functionality for verifier-side Fiat_Shamir challenge generation.
use ark_bn254::Fr;
use ark_ff::{BigInteger, PrimeField};
use sha2::{Digest, Sha256};

use crate::error::GnarkError;

/// A simple wrapper that accumulates bytes and hashes them to a field element.
pub(crate) struct WrappedHashToField {
    domain: Vec<u8>,
    to_hash: Vec<u8>,
}

impl WrappedHashToField {
    /// Create a new instance with a given domain separator.
    pub(crate) fn new(domain_separator: &[u8]) -> Self {
        Self {
            domain: domain_separator.to_vec(),
            to_hash: Vec::new(),
        }
    }

    /// Append bytes to the internal buffer.
    pub fn write(&mut self, data: &[u8]) {
        self.to_hash.extend_from_slice(data);
    }

    /// Hash accumulated bytes to a field element (returning its byte representation).
    pub fn sum(&mut self, mut b: Vec<u8>) -> Vec<u8> {
        let res = hash_to_field(&self.to_hash, &self.domain, 1)
            .expect("Unable to get elements from hashing");

        // Convert the hash output to a field element
        let element = res[0];

        // Clear the buffer (optional; similar to Reset)
        self.to_hash.clear();

        // Return the field element as bytes (big-endian)
        b.extend_from_slice(&element.into_bigint().to_bytes_be());

        b
    }

    /// Reset internal state (clears accumulated bytes)
    pub fn reset(&mut self) {
        self.to_hash.clear();
    }
}

/// Hashes to field elements with 128-bit security.
/// Equivalent to the Go `Hash` function.
pub(crate) fn hash_to_field(msg: &[u8], dst: &[u8], count: usize) -> Result<Vec<Fr>, GnarkError> {
    // 128 bits of security
    // L = ceil((ceil(log2(p)) + k) / 8), where k = 128
    let bits = Fr::MODULUS_BIT_SIZE as usize;
    let bytes = 1 + (bits - 1) / 8;
    let l = 16 + bytes; // per RFC 9380
    let len_in_bytes = count * l;

    // Step 1: Expand message
    let pseudo_random_bytes = expand_message_xmd(msg, dst, len_in_bytes)?;

    // Step 2: Convert to field elements
    let mut res = Vec::with_capacity(count);
    for i in 0..count {
        let start = i * l;
        let end = (i + 1) * l;
        let slice = &pseudo_random_bytes[start..end];

        // Interpret as a big integer mod p
        let fe = Fr::from_be_bytes_mod_order(slice);
        res.push(fe);
    }

    Ok(res)
}

/// ExpandMsgXmd per RFC 9380 §5.4.1 using SHA-256.
/// Expands `msg` and `dst` into `len_in_bytes` pseudorandom bytes.
fn expand_message_xmd(msg: &[u8], dst: &[u8], len_in_bytes: usize) -> Result<Vec<u8>, GnarkError> {
    const B_IN_BYTES: usize = 64; // SHA-256 block size (in bytes)
    const H_LEN: usize = 32; // SHA-256 output size (in bytes)

    let ell = len_in_bytes.div_ceil(H_LEN);
    if ell > 255 {
        return Err(GnarkError::HashError(
            "invalid len_in_bytes: too large".into(),
        ));
    }
    if dst.len() > 255 {
        return Err(GnarkError::HashError(
            "invalid DST length (>255 bytes)".into(),
        ));
    }

    let size_domain = dst.len() as u8;
    let mut dst_prime = Vec::with_capacity(dst.len() + 1);
    dst_prime.extend_from_slice(dst);
    dst_prime.push(size_domain);

    // Step 1: b0 = H(Z_pad || msg || I2OSP(len_in_bytes, 2) || I2OSP(0, 1) || DST_prime)
    let mut hasher = Sha256::new();
    hasher.update(vec![0u8; B_IN_BYTES]); // Z_pad
    hasher.update(msg);
    hasher.update([(len_in_bytes >> 8) as u8, (len_in_bytes & 0xff) as u8]);
    hasher.update([0x00]);
    hasher.update(&dst_prime);
    let b0 = hasher.finalize();

    // Step 2: b1 = H(b0 || I2OSP(1, 1) || DST_prime)
    let mut hasher = Sha256::new();
    hasher.update(b0);
    hasher.update([0x01]);
    hasher.update(&dst_prime);
    let mut bi = hasher.finalize();

    // Step 3: accumulate blocks
    let mut result = Vec::with_capacity(len_in_bytes);
    result.extend_from_slice(&bi);

    for i in 2..=ell {
        // bi = H(strxor(b0, b_{i-1}) || I2OSP(i,1) || DST_prime)
        let mut strxor = [0u8; H_LEN];
        for j in 0..H_LEN {
            strxor[j] = b0[j] ^ bi[j];
        }
        let mut hasher = Sha256::new();
        hasher.update(strxor);
        hasher.update([i as u8]);
        hasher.update(&dst_prime);
        bi = hasher.finalize();

        let start = (i - 1) * H_LEN;
        let end = std::cmp::min(i * H_LEN, len_in_bytes);
        result.extend_from_slice(&bi[..end - start]);
    }

    result.truncate(len_in_bytes);
    Ok(result)
}
