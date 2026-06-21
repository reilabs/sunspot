---
title: Installation
nav_order: 2
---

# Installation

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Solana tools](https://solana.com/docs/intro/installation)
- Noir **1.0.0-beta.22**

## Build from source

```bash
git clone git@github.com:reilabs/sunspot.git
cd sunspot/go
go build -o sunspot .
```

## Add `sunspot` to your `PATH`

The easiest option is to move the binary into a directory already on your `PATH`:

```bash
sudo mv sunspot /usr/local/bin/
```

Or create a personal `bin` folder:

```bash
mkdir -p ~/bin
mv sunspot ~/bin/
```

Then add it to your shell config:

**bash**
```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bash_profile
source ~/.bash_profile
```

**zsh** (default on macOS)
```bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

Verify the install:

```bash
sunspot --help
```

## Configure `GNARK_VERIFIER_BIN`

`sunspot deploy` needs `GNARK_VERIFIER_BIN` to point at the `verifier-bin` crate directory:

```bash
export GNARK_VERIFIER_BIN=/path/to/verifier-bin
```

Add this line to:

- **bash (Linux):** `~/.bashrc`
- **bash (macOS):** `~/.bash_profile`
- **zsh:** `~/.zshrc`

Reload your shell so the variable is picked up:

```bash
source ~/.zshrc   # or ~/.bashrc, ~/.bash_profile
```
