# Sanction Checker example
This is an example of we can use [Sunspot](https://github.com/reilabs/sunspot/) to verify noir circuits on solana. 
A token is minted to an account only if a proof of valid non-ofac sanctioned passport is provided.

> Note: the circuit is designed only for passports that use RSA for both the CSC and the DSC, using PKCS signatures for the CSC and PSS signatures for the DSC.

## Project Setup
The circuit is defined in the `sanction_checker_circuit` folder. 
The verification key is derived by compiling the Noir project with `nargo compile` and using the sunspot project to generate the Gnark verification and proving keys.
Proof and public witness generation is performed by submitting valid non-sanctioned passport inputs to the circuit and calling `nargo execute`, then using Sunspot to generate the Gnark proof and public witness.

The program that mints the token upon submission of a valid proof and public witness is defined in the `program` directory.

A client to interact with the deployed minting program is defined in the `client` directory.

The verification key used and an example proof and public witness, can be found in the `test_files` directory.

## Statistics
Invoking the verification circuit incurs costs of 471486 compute units. 
The proof size is 388 bytes. 