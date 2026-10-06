---
title: Home
layout: home
nav_order: 1
---

# Sunspot
{: .fs-9 }

Prove and verify [Noir](https://noir-lang.org) circuits using Groth16.
{: .fs-6 .fw-300 }

[Get started](./installation){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/reilabs/sunspot){: .btn .fs-5 .mb-4 .mb-md-0 }

---

{: .note }
Sunspot has been audited by [HashCloak](https://hashcloak.com). See the [audit report](https://github.com/reilabs/sunspot/blob/main/audits/2026_09_hashcloak_audit_report.pdf).
Requires **Noir 1.0.0-beta.22**.

## What is Sunspot?

Sunspot is a toolchain that bridges [Noir](https://noir-lang.org) — a zero-knowledge proof DSL —
with the [Solana](https://solana.com) blockchain. It lets you:

- **Compile** Noir ACIR circuits into gnark-compatible constraint systems (CCS).
- **Generate** Groth16 proving and verifying keys.
- **Prove and verify** circuit executions off-chain.
- **Deploy** a verifying Solana program that can check those proofs on-chain.

## Quickstart

```bash
# Compile a Noir ACIR file into a CCS file
sunspot compile my_circuit.json

# Generate a proving and verifying key
sunspot setup my_circuit.ccs

# Create a Groth16 proof
sunspot prove my_circuit.json witness.gz my_circuit.ccs proving_key.pk

# Verify a proof locally
sunspot verify verifying_key.vk proof.proof public_witness.pw

# Deploy a Solana verifier program for your circuit
sunspot deploy verifying_key.vk
```

See the [Commands reference](./commands) for full details.

