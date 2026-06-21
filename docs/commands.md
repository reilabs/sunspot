---
title: Commands
nav_order: 3
---

# CLI Reference

All commands follow the form:

```bash
sunspot [command] [args...]
```

For per-command help:

```bash
sunspot [command] --help
```

## Overview

| Command      | Description                                                                                    |
| ------------ | ---------------------------------------------------------------------------------------------- |
| `compile`    | Compile an ACIR file into a CCS file.                                                          |
| `setup`      | Generate a proving key (pk) and verifying key (vk) from a CCS file.                            |
| `prove`      | Generate a Groth16 proof and public witness from an ACIR file, witness, CCS, and proving key. |
| `verify`     | Verify a proof and public witness with a verification key.                                     |
| `deploy`     | Create a verifying Solana program executable and keypair.                                      |
| `completion` | Generate the shell autocompletion script.                                                      |
| `help`       | Display help for any command.                                                                  |

## Examples

### Compile

```bash
sunspot compile my_circuit.json
```

### Setup

{: .warning }
`setup` performs a gnark trusted setup with **no mitigation for cryptographic toxic waste**.

```bash
sunspot setup my_circuit.ccs
```

### Prove

```bash
sunspot prove my_circuit.json witness.gz my_circuit.ccs proving_key.pk
```

### Verify

```bash
sunspot verify verifying_key.vk proof.proof public_witness.pw
```

### Deploy

```bash
sunspot deploy verifying_key.vk
```

Requires [`GNARK_VERIFIER_BIN`](./installation#configure-gnark_verifier_bin) environment variable to be set and pointing to the `verifier-bin` crate from this project.
