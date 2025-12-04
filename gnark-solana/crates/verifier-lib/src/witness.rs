//! Provides utilities for parsing Gnark-generated public witnesses
use std::io::{self, Read};

use crate::error::GnarkError;

/// The Gnark witness — public inputs to the circuit.
pub struct GnarkWitness<const NR_INPUTS: usize> {
    /// The variables in the public witness
    pub entries: [[u8; 32]; NR_INPUTS],
}

impl<const NR_INPUTS: usize> GnarkWitness<NR_INPUTS> {
    /// Parses the witness (public inputs) from a reader.
    pub fn parse<R: Read>(mut reader: R) -> io::Result<Self> {
        // We first parse through 12 bytes,
        // Gnark witness encoding encodes the number of public inputs in 4 bytes,
        // and the number of private inputs in 4 bytes, even though in the public witness
        // the number of private inputs is always zero.
        // Then the underlying vector type encodes the number of entries in 4 bytes
        let mut len_buf = [0u8; 12];
        reader.read_exact(&mut len_buf)?;

        let mut entries = [[0u8; 32]; NR_INPUTS];
        for entry in entries.iter_mut() {
            reader.read_exact(entry)?;
        }
        Ok(Self { entries })
    }

    /// Constructs a witness directly from a byte slice.
    /// Expects the same layout as `parse()`: 12-byte header + NR_INPUTS * 32 bytes of data.
    pub fn from_bytes(bytes: &[u8]) -> Result<Self, GnarkError> {
        const HEADER_LEN: usize = 12;
        let expected_len = HEADER_LEN + NR_INPUTS * 32;

        if bytes.len() != expected_len {
            return Err(GnarkError::PublicWitnessParsingError);
        }

        // Skip 12-byte header
        let mut entries = [[0u8; 32]; NR_INPUTS];
        let mut offset = HEADER_LEN;

        for entry in entries.iter_mut() {
            entry.copy_from_slice(&bytes[offset..offset + 32]);
            offset += 32;
        }

        Ok(Self { entries })
    }
}

#[cfg(test)]
mod tests {
    use std::fs::File;

    use num_bigint::BigUint;
    use num_traits::Num;
    #[test]
    fn test_parse_witness_from_file() {
        let mut file = File::open("src/test_files/sum_a_b.pw").unwrap();

        const NR_INPUTS: usize = 1;
        let witness = super::GnarkWitness::<NR_INPUTS>::parse(&file);

        assert!(witness.is_ok());
        let witness = witness.unwrap();

        // Assert that we have hit the end of the witness file by trying to read one more byte
        let mut buf = [0u8; 1];
        let bytes_read = std::io::Read::read(&mut file, &mut buf).unwrap();
        assert_eq!(
            bytes_read, 0,
            "Expected EOF after parsing witness, but more bytes remain!"
        );

        let value = BigUint::from_bytes_be(&witness.entries[0]);
        let expected = BigUint::from_str_radix("5000", 10).unwrap();
        assert_eq!(value, expected);
    }
}
