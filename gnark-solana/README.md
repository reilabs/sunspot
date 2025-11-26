# Gnark - Solana

This project provides tools to verify [Gnark](https://github.com/Consensys/gnark) generated groth16 proofs on [Solana](https://solana.com/).

## Prerequisites

Install Solana CLI and other requisite tools: 

```bash
curl --proto '=https' --tlsv1.2 -sSfL https://solana-install.solana.workers.dev | bash
```

Confirm the installation:

```bash
solana --version
rustc --version
```

Optionally, configure your Solana CLI to use the **devnet**:

```bash
solana config set --url https://api.devnet.solana.com
```

## Compile Solana Program

1. Generate a verification key with [Sunspot](https://github.com/reilabs/sunspot).

2. Place the newly generated verification key file inside the `verifier-bin` folder of this project.

   * For example: `verifier-bin/my_verification_key.vk`

Then run the following command:

```bash
VK_PATH="my_verification_key.vk" cargo build-sbf
```

> **Note:** If provided a relative path, the build process automatically looks for the verification key inside the `verifier-bin` folder.
> If the `VK_PATH` environment variable is not set, it will default to using:
>
> ```
> verifier-bin/default.vk
> ```

This will generate the Solana program and its keypair files in the `target/deploy` folder.

## Deploy Solana Program

Once the program is compiled, you can deploy it to your desired Solana cluster.

1. **Choose your cluster** (devnet, testnet, or mainnet):

```bash
# Example: set cluster to devnet
solana config set --url https://api.devnet.solana.com
```

2. **Deploy the program** using the compiled `verifier-bin.so` file in `target/deploy`:

```bash
solana program deploy target/deploy/verifier-bin.so
```

* The deployment command will return a **program ID**, which you will use to interact with the program.

3. **Verify the deployment** (optional):

```bash
solana program show <program_id>
```

This confirms that your program is deployed and ready to use on the selected Solana cluster.

## Statistics

Verification of a proof and public witness will cost between 170000 and 500000 compute units, depending on the complexity of the circuit. 
Proof sizes will be between 324 and 388 bytes.