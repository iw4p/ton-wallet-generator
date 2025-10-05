# TON Wallet Manager

A comprehensive **offline** Go tool for managing TON wallets. Generate wallets from seed phrases, manage multiple wallet versions, and handle subwallet IDs with a clean, professional interface.

## Features

- **Completely offline** - No network connection required
- **Interactive & CLI modes** - User-friendly interface with scripting support
- **Multiple wallet versions** - Support for V3R1, V3R2, V4R1, V4R2, V5R1 Final
- **Network selection** - Mainnet/Testnet for proper address generation
- **Subwallet management** - Create multiple wallets from one seed
- **Flexible input** - Accept seed phrases in any format (spaces, commas, lines)
- **Clean output** - Parseable CLI output for automation
- **Secure derivation** - TON standard PBKDF2 + HMAC-SHA512
- **Professional structure** - Clean, maintainable Go architecture

## Installation

### Quick Install

```bash
# Build the application
go build -o bin/ton-wallet-manager ./cmd/ton-wallet-manager

# Or install to $GOPATH/bin
go install ./cmd/ton-wallet-manager
```

### Manual Build

```bash
go build -o bin/ton-wallet-manager ./cmd/ton-wallet-manager
```

The binary will be created in the `bin/` directory.

## Quick Start

```bash
# Interactive mode (recommended for first-time users)
./bin/ton-wallet-manager

# CLI mode (for automation and scripting)
./bin/ton-wallet-manager --generate --simple
```

## Usage

### Interactive Mode

Run the program without any flags for interactive mode:

```bash
./bin/ton-wallet-manager
# Or if installed: ton-wallet-manager
```

The program will interactively ask you:

1. **Generate new seed or use existing**: Choose whether to create a new seed phrase or input an existing one
   - **Flexible input format**: When entering an existing seed, you can paste it in any format:
     - Single line with spaces: `word1 word2 word3...`
     - Comma-separated: `word1, word2, word3...`
     - Multiple lines (press Enter twice when done)
     - Mixed formats with tabs, semicolons, or other separators
2. **Network selection**: Choose Mainnet or Testnet
3. **Wallet version**: Select from V3R1, V3R2, V4R1, V4R2, or V5R1 Final (V4R2 recommended)
4. **Subwallet ID**: Enter a custom subwallet ID or press Enter for default (698983191)

### CLI Mode (Non-Interactive)

Use CLI flags for non-interactive wallet generation. Clean output suitable for piping and scripting:

```bash
# Generate a new wallet and show only the address (cleanest output)
./bin/ton-wallet-manager --generate --simple

# Generate a new wallet with parseable output
./bin/ton-wallet-manager --generate

# Generate a wallet on testnet with V5R1
./bin/ton-wallet-manager --generate --network testnet --version v5r1

# Use an existing seed phrase
./bin/ton-wallet-manager --seed "word1 word2 word3 ... word24" --simple

# Create a wallet with custom subwallet ID
./bin/ton-wallet-manager --generate --subwallet 0

# Pipe examples - extract specific fields
./bin/ton-wallet-manager --generate | grep "^Address:"
./bin/ton-wallet-manager --generate | awk -F': ' '/^Address:/ {print $2}'
./bin/ton-wallet-manager --generate --simple | xargs -I {} echo "New wallet: {}"
```

#### Available CLI Flags

- `--generate` - Generate a new seed phrase
- `--seed <phrase>` - Use an existing 24-word seed phrase (space-separated)
- `--network <name>` - Network: `mainnet` or `testnet` (default: mainnet)
- `--version <ver>` - Wallet version: `v3r1`, `v3r2`, `v4r1`, `v4r2`, `v5r1beta`, `v5r1final` (default: v4r2)
- `--subwallet <id>` - Custom subwallet ID (default: 698983191 for v3/v4, 0 for v5)
- `--simple` - Simple output mode (only show the address on stdout)

**Note:** Either `--generate` or `--seed` must be provided in CLI mode.

**About `--simple` flag:**

- With `--generate`: Seed phrase goes to **stderr** (visible but not piped), address goes to **stdout** (pipeable)
- With `--seed`: Only address goes to **stdout** (useful when you already have the seed and just need the address)

**Default Settings:**
When using `--generate --simple` (or any CLI flags without specifying options), the tool uses these defaults:

- **Network**: Mainnet
- **Version**: V4R2 (recommended)
- **Subwallet ID**: 698983191 (standard for V4R2)
- **Output**: Address only to stdout, seed phrase to stderr

**Example:**

```bash
# This creates a Mainnet V4R2 wallet with subwallet 698983191
./bin/ton-wallet-manager --generate --simple
```

#### CLI Output Format

**Simple mode with existing seed** (`--seed "..." --simple`):

```
EQAwSf6OBbhbNNTaJWTbH_rxIPpi-JYuZz67qnJVXDEuDfn5
```

**Simple mode with generation** (`--generate --simple`):

```
# Seed phrase (save securely!):           ← stderr (visible in terminal)
word1 word2 ... word24                    ← stderr (visible in terminal)
EQAwSf6OBbhbNNTaJWTbH_rxIPpi-JYuZz67qnJVXDEuDfn5  ← stdout (pipeable)
```

**Regular CLI mode** (clean, parseable):

```
Address: EQAwSf6OBbhbNNTaJWTbH_rxIPpi-JYuZz67qnJVXDEuDfn5
Network: Mainnet
Version: V4R2
Subwallet: 698983191
PublicKey: 4dab400f4b6e682ee706472973b1bf49c3ef09c9c4fc79ae86745f105c8d418a
PrivateKey: dcbe502c265c460dd3d1c7b6fafa694741f930108624a12b2e8dc72bbf3dc9ab4dab400f4b6e682ee706472973b1bf49c3ef09c9c4fc79ae86745f105c8d418a
```

**Important:** When generating with `--generate`, the seed is **always** output (to stderr in `--simple` mode, to stdout otherwise). This ensures you never lose access to your wallet.

#### Practical Examples

```bash
# Generate a wallet and save address to variable (seed visible in terminal)
ADDRESS=$(./bin/ton-wallet-manager --generate --simple 2>&1 | tail -1)

# Get address from existing seed for scripting
./bin/ton-wallet-manager --seed "your 24 words here" --simple

# Generate wallet and filter output for parsing
./bin/ton-wallet-manager --generate | grep "^Address:" | cut -d' ' -f2

# Create multiple wallets with different subwallet IDs from same seed
for i in 0 1 2; do
  ./bin/ton-wallet-manager --seed "your seed" --subwallet $i --simple
done
```

### Creating Multiple Wallets

To create multiple wallets from the same seed phrase, run the program multiple times with:

- Same seed phrase
- Different subwallet IDs (e.g., 0, 1, 2, 3...)
- Same or different wallet versions

### Example

```
TON Wallet Manager
===================

Generate new seed phrase? (y/n): y

Your new seed phrase (SAVE THIS SECURELY):
word1 word2 word3 ... word24

Select network:
1. Mainnet
2. Testnet

Select network (1-2): 1
Using Mainnet

Available wallet versions:
1. V3R1
2. V3R2
3. V4R1
4. V4R2 (recommended)
5. V5R1 Final

Select wallet version (1-5): 4

Enter subwallet ID (press Enter for default 698983191): 0

========================================
Wallet Created Successfully!
========================================
Network: Mainnet
Version: V4R2
Subwallet ID: 0
Address: EQ...
Public Key: abc123...
Private Key: def456...
========================================

Note: Store your seed phrase securely!
This is an offline wallet manager - no network connection required!
```

## Security Notes

- **Never share your seed phrase** - Anyone with access to it can control all wallets derived from it
- **Store your seed phrase securely** - Write it down and keep it in a safe place
- **Private keys are sensitive** - The program displays private keys for convenience, but never share them
- Different subwallet IDs create completely different wallet addresses from the same seed

## Project Structure

This project is organized into clean, maintainable packages:

```
ton-wallet-manager/
├── cmd/ton-wallet-manager/    # Application entry point
├── internal/
│   ├── wallet/                # Core wallet management logic
│   ├── cli/                   # Command-line interface
│   └── ui/                    # User interface (interactive & display)
├── bin/                       # Build output
└── STRUCTURE.md               # Detailed architecture documentation
```

For detailed information about the architecture and design decisions, see [STRUCTURE.md](STRUCTURE.md).

## Technical Details

- **Completely offline** - No network dependencies
- Uses PBKDF2 with HMAC-SHA512 for key derivation (TON standard)
- Follows TON's standard mnemonic-to-key conversion
- Supports all major wallet contract versions (V3R1, V3R2, V4R1, V4R2, V5R1 Final)
- Network-aware address generation (Mainnet/Testnet)
- Default subwallet ID: 698983191
- Secure key derivation using TON's official method
- Clean architecture with separation of concerns
- Factory functions over classes (functional approach)
- Professional Go project structure
- Easy to extend and maintain
- Simple build process with standard Go tools

## ⚠️ Disclaimer

**USE AT YOUR OWN RISK**

This software is provided "as is" without warranty of any kind. The author is not responsible for:

- **Loss of funds** due to incorrect usage, lost seed phrases, or private keys
- **Security vulnerabilities** or bugs in the software
- **Compatibility issues** with different TON wallet implementations
- **Any financial losses** resulting from the use of this tool

**Important Security Notes:**

- Always verify addresses on multiple sources before sending funds
- Store your seed phrases securely offline
- Test with small amounts first
- This tool is for educational and personal use only
- The author assumes no responsibility for any losses or damages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
