# Sunspot

**Sunspot** provides tools to prove and verify [noir](https://noir-lang.org) circuits on solana.

> ⚠️ Requires **Noir 1.0.0-beta.18**

## Security

**Sunspot has not been audited yet and is provided as-is. We make no guarantees to its safety or reliability.**

To report security vulnerabilities, please use the `Security` tab on this repository.

## Installation

Make sure you have [Go 1.24+](https://go.dev/dl/) and [Solana tools](https://solana.com/docs/intro/installation) installed.

```bash
# Clone the repository
git clone git@github.com:reilabs/sunspot.git
cd sunspot/go

# Build the binary
go build -o sunspot .
````

#### Add the binary to your PATH

You can move the binary to a directory already in your `PATH` (easiest):

```bash
sudo mv sunspot /usr/local/bin/
```

Alternatively, you can create a `bin` folder in your home directory and add it to your PATH.

```bash
# Create a personal bin folder if you don’t have one
mkdir -p ~/bin
mv sunspot ~/bin/
```

Then add this line to your shell configuration file:

* For **bash**:

  ```bash
  echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bash_profile
  source ~/.bash_profile
  ```

* For **zsh** (default on macOS):

  ```bash
  echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
  source ~/.zshrc
  ```

Now you can run `sunspot` from anywhere:

```bash
sunspot --help
```

## Setting `GNARK_VERIFIER_BIN`

`GNARK_VERIFIER_BIN` must point to the verifier-bin crate directory in order for `sunspot deploy` to work.

To configure this, add the following line to your shell’s configuration file:

```bash
export GNARK_VERIFIER_BIN=/path/to/verifier-bin
```

### Where to add it

* **bash (Linux):** add the line to `~/.bashrc`
* **bash (macOS):** add the line to `~/.bash_profile`
* **zsh:** add the line to `~/.zshrc`

After editing the file, reload your shell:

```bash
source ~/.bashrc       # or ~/.bash_profile, ~/.zshrc depending on your shell
```

Now `GNARK_VERIFIER_BIN` will be available in all future terminal sessions.


## Usage

After installing **Sunspot**, you can use it as a command-line tool for working with **Noir circuits on Solana**.

```bash
sunspot [command]
````


###  Available Commands

| Command      | Description                                                                      |
| ------------ | -------------------------------------------------------------------------------- |
| `compile`    | Compile an ACIR file into a CCS file                                             |
| `completion` | Generate the autocompletion script for the specified shell                       |
| `help`       | Display help information about any command                                       |
| `prove`      | Generate a Groth16 proof and public witness from an ACIR file, a witness, CCS, and proving key |
| `setup`      | Generate a proving key (pk) and verifying key (vk) from a CCS file               |
| `verify`     | Verify a proof and public witness with a verification key                        |
| `deploy`     | Create a verifying solana program executable and keypair|

### 💡 Examples

```bash
# Compile a Noir ACIR file into a CCS file
sunspot compile my_circuit.json

# Generate a proving and verifying key
# ⚠️ THIS IS UNSAFE!
# ⚠️ IT PERFORMS GNARK TRUSTED SETUP WITH NO MITIGATION FOR CRYPTOGRAPHIC TOXIC WATE!
sunspot setup my_circuit.ccs

# Create a Groth16 proof
sunspot prove my_circuit.json witness.gz my_circuit.ccs proving_key.pk

# Verify a proof
sunspot verify verifying_key.vk proof.proof public_witness.pw

# Create Solana verification program
sunspot deploy verifying_key.vk 
```

For detailed information on each command:

```bash
sunspot [command] --help
```

## Codebase Overview

This project is organized as follows:

- `go/` – Contains functionality to parse Noir circuits and witnesses and produces gnark outputs, also contains CLI functionality in `go/cmd` subdirectory.
- `gnark-solana/` – Provides functionality to verify gnark proofs on solana, a fuller description of this directory can be found [here](gnark-solana/README.md).
- `noir-samples/` – Example Noir projects used for unit and integration tests.
- `testgen` - Creates ACIR snippets to test parsing, does **not** produce semantically valid programs.


## Credits

- **Light Protocol**  
 Our thanks goes to Light protocol, the original authors of the [Groth16-solana](https://github.com/Lightprotocol/groth16-solana) repo, who published it under the Apache 2.0 License.
 We used this for inspiration for our own core Gnark verifier for both the error type definition and  some of the core verifier functionality.
