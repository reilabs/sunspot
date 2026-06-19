# Installing Sunspot

Sunspot can be installed either from a prebuilt release artifact or from source. In both cases, you will also need to point Sunspot to the `verifier-bin` crate if you intend to deploy on Solana.

## Prerequisites

- [Solana tools](https://solana.com/docs/intro/installation) — required for deployment on Solana.
- [Go 1.24+](https://go.dev/dl/) — required only if you are building from source.

## Install from a release

Prebuilt artifacts are published on the [Sunspot releases page](https://github.com/reilabs/sunspot/releases). Each release contains:

- `sunspot_<version>_<os>_<arch>.tar.gz` — archives for Linux and macOS on `amd64` and `arm64` architectures.
- `sunspot_<version>_<arch>.deb` — Debian/Ubuntu packages.
- `sunspot_<version>_<arch>.rpm` — Fedora/RHEL packages.
- `checksums.txt` — SHA-256 checksums for every artifact in the release.


### macOS / Linux (tar.gz)

1. Download the archive matching your OS and architecture from the releases page.
2. Extract it and move the binary somewhere on your `PATH`:

   ```bash
   tar -xzf sunspot_<version>_<os>_<arch>.tar.gz
   sudo mv sunspot /usr/local/bin/
   ```

3. On macOS, the binary is unsigned. The first time you run it you may need to allow it under **System Settings → Privacy & Security**, or remove the quarantine attribute:

   ```bash
   xattr -d com.apple.quarantine /usr/local/bin/sunspot
   ```

### Debian / Ubuntu (.deb)

```bash
sudo dpkg -i sunspot_<version>_<arch>.deb
```

This installs the binary to `/usr/bin/sunspot`.

### Fedora / RHEL (.rpm)

```bash
sudo rpm -i sunspot_<version>_<arch>.rpm
```

This installs the binary to `/usr/bin/sunspot`.

### Verify the installation

```bash
sunspot --help
```

## Install from source

```bash
# Clone the repository
git clone git@github.com:reilabs/sunspot.git
cd sunspot/go

# Build the binary
go build -o sunspot ./cmd/sunspot
```

### Add the binary to your PATH

You can move the binary to a directory already in your `PATH` (easiest):

```bash
sudo mv sunspot /usr/local/bin/
```

Alternatively, create a `bin` folder in your home directory and add it to your `PATH`:

```bash
mkdir -p ~/bin
mv sunspot ~/bin/
```

Then add this line to your shell configuration file:

- For **bash**:

  ```bash
  echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bash_profile
  source ~/.bash_profile
  ```

- For **zsh** (default on macOS):

  ```bash
  echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
  source ~/.zshrc
  ```

Now you can run `sunspot` from anywhere:

```bash
sunspot --help
```

## Configure `GNARK_VERIFIER_BIN`

`GNARK_VERIFIER_BIN` must point to the `verifier-bin` crate directory in order for `sunspot deploy` to work. The release artifacts only ship the `sunspot` binary, so you will need a clone of this repository for the verifier crate regardless of how you installed Sunspot:

```bash
git clone git@github.com:reilabs/sunspot.git
export GNARK_VERIFIER_BIN=/path/to/sunspot/gnark-solana/crates/verifier-bin
```

To persist the setting, add the `export` line to your shell configuration file:

- **bash (Linux):** `~/.bashrc`
- **bash (macOS):** `~/.bash_profile`
- **zsh:** `~/.zshrc`

After editing, reload your shell:

```bash
source ~/.bashrc       # or ~/.bash_profile, ~/.zshrc depending on your shell
```

`GNARK_VERIFIER_BIN` will now be available in all future terminal sessions.
