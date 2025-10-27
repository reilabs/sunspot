# Sunspot

**Sunpot** provides tools to prove and verify [noir](https://noir-lang.org) circuits on solana.

> 🚧 **Work in Progress:** Sunspot as is provides a Groth16 backend to ACIR circuit representations. It does not yet support verifying the generated proofs on Solana.

## Installation

Make sure you have [Go 1.24+](https://go.dev/dl/) installed.

```bash
# Clone the repository
git clone git@github.com:reilabs/sunspot.git
cd sunspot

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

### 💡 Examples

```bash
# Compile a Noir ACIR file into a CCS file
sunspot compile my_circuit.json

# Generate a proving and verifying key
sunspot setup my_circuit.ccs

# Create a Groth16 proof
sunspot prove my_circuit.json witness.gz my_circuit.ccs proving_key.pk

# Verify a proof
sunspot verify verifying_key.vk proof.proof public_witness.pw
```

For detailed information on each command:

```bash
sunspot [command] --help
```