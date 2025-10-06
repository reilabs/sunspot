# `gnark` integration for Noir

This project provides a Noir interface for the gnark library, allowing Noir programs to utilize zk-SNARKs for zero-knowledge proofs. It includes implementations of various cryptographic functions and primitives, enabling developers to create privacy-preserving applications.


# Black box functions implemented

|Function|Implemented|
|--------|-----------|
|AES128Encrypt|Yes|
|AND|Yes|
|XOR|Yes|
|Range|Yes|
|Blake2s|Yes|
|Blake3|Yes|
|ECDSA SECP256K1|Yes|
|ECDSA SECP256R1|Yes|
|MultiScalarMul|Yes|
|KeccakF1600|Yes|
|RecursiveAggregation|No|
|EmbeddedCurveAdd|Yes|
|Poseidon2Permutation|Yes|
|SHA256Compression|Yes|