# gnark integration for Noir

This project provides a Noir interface for the gnark library, allowing Noir programs to utilize zk-SNARKs for zero-knowledge proofs. It includes implementations of various cryptographic functions and primitives, enabling developers to create privacy-preserving applications.


# Black box functions implemented

|Function|Implemented|
|--------|-----------|
|AES128Encrypt|Yes|
|AND|Yes|
|BigIntAdd|Yes|
|BigIntDiv|Yes|
|BigIntFromLEBytes|No|
|BigIntMul|Yes|
|BigIntSub|Yes|
|BigIntToLEBytes|No|
|Blake2s|Yes|
|Blake3|Yes|
|ECDSA SECP256K1|Yes|
|ECDSA SECP256R1|Yes|
|Keccak F1600|No|
|MultiScalarMul|No|
|Poseidon2Permutation|WIP|
|Range|Yes|
|RecursiveAggregation|No|
|SHA256Compression|Yes|
|XOR|Yes|