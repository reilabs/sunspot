# ACIR support in gnark

This document contains the information about the way ACIR circuit support is maintained in ```gnark```.

# Noir/ACIR vs gnark - key differences 

As discussed in details in the [ACIR format document](acir.md) - ACIR format is a generic ZK-circuit format produced by Noir language compiler. ACIR format is generic and can support any custom fields, without being limited to BN254 elliptic curve field or the specific set of fields. Therefore, in order to have the ACIR format decoding properly supported, ACIR structures in Go representations use the ```ACIRField``` interface generic requiring the routines for decoder and comparison of two separate fields.

```gnark``` circuits using the circuit API, however, are typically able to natively use the limited set of curves implemented in [```github.com/consensys/gnark-crypto/ecc```](https://github.com/Consensys/gnark-crypto/tree/master/ecc) package. This list includes the following curves:
 * [BLS12_377](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bls12-377)
 * [BLS12_381](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bls12-381)
 * [BLS24_315](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bls24-315)
 * [BLS24_317](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bls24-317)
 * [BN254](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bn254)
 * [BW6_633](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bw6-633)
 * [BW6_761](https://github.com/Consensys/gnark-crypto/tree/master/ecc/bw6-761)
 * [Grumpkin](https://github.com/Consensys/gnark-crypto/tree/master/ecc/grumpkin)
 * [SECP256k1](https://github.com/Consensys/gnark-crypto/tree/master/ecc/secp256k1)
 * [STARK curve](https://github.com/Consensys/gnark-crypto/tree/master/ecc/stark-curve)

# Solution

As it can be seen, before building (defining, compiling) the ACIR circuit into R1CS using Circuit API we need to convert the ACIR fields into gnark-compatible ones. This is achieved in our solution through the following stages

## Decoding the ACIR

Decoding the binary code is performed into ACIR Decoded structures using the field template of ```ACIRField``` interface. This interface contains the following routines:
 * ```UnmarshalReader(r io.Reader) error``` - decoding the binary input (from file reader or some other gzip7 decompressor) into the structure
 * ```Equals(other ACIRField)``` - comparison routine, useful for testing purposes - allows comparing ACIR fields directly
 * ```ToGenericFpElement``` - converting the ACIR field into the generic FP element used as the generic container for supported fields in ```gnark```

Decoded structures have the ```Decoded``` annotation (e.g. ```ExpressionDecoded[T ACIRField]```)

## Converting the decoded structures into gnark-friendly ones


