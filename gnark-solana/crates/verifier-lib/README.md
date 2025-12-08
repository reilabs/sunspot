# Gnark Verifier Lib

This crate provides utilities for verifying [Gnark](https://github.com/Consensys/gnark) proofs on Solana.

The verifier expects all constructs. i.e verification keys, proofs and public witnesses to be serialized in raw, non-compressed, format. This lowers on-chain compute costs by skipping the decompression overhead.

Use the `vk` file to generate a compile-time constant verification key for your verifier.

Verifier costs 175,125 compute units to invoke without any commitments and 471486 to invoke with 1 commitment.