# `gnark` integration for Noir

This project provides a Noir interface for the gnark library, allowing Noir programs to utilize zk-SNARKs for zero-knowledge proofs. It includes implementations of various cryptographic functions and primitives, enabling developers to create privacy-preserving applications.


# Black box functions implemented

|Function|Implemented|
|--------|-----------|
|AES128Encrypt|No|
|AND|Yes|
|XOR|Yes|
|Range|Yes|
|SHA256Hash|No|
|Blake2s|No|
|Blake3|No|
|SchnorrVerify|No|
|PedersenCommitment|No|
|PedersenHash|No|
|ECDSA SECP256K1|No|
|ECDSA SECP256R1|No|
|MultiScalarMul|No|
|Keccak256|No|
|KeccakF1600|Yes|
|RecursiveAggregation|No|
|EmbeddedCurveAdd|No|
|Poseidon2Permutation|No|
|SHA256Compression|Yes|