# TON Wallet Generator

A simple **offline** Go tool to generate TON wallets from seed phrases using different wallet versions and subwallet IDs.

## Features

- **Completely offline** - No network connection required
- Generate new 24-word seed phrases or use existing ones
- Support for multiple wallet versions (V3R1, V3R2, V4R1, V4R2, V5R1 Final)
- Network selection (Mainnet/Testnet) for proper address generation
- Custom subwallet ID support for creating multiple wallets from one seed
- Displays wallet address, public key, and private key
- Secure key derivation using TON's standard method

## Installation

```bash
go build -o ton-wallet-generator
```

## Usage

### Interactive Mode

Run the program without any flags for interactive mode:

```bash
./ton-wallet-generator
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
./ton-wallet-generator --generate --simple

# Generate a new wallet with parseable output
./ton-wallet-generator --generate

# Generate a wallet on testnet with V5R1
./ton-wallet-generator --generate --network testnet --version v5r1

# Use an existing seed phrase
./ton-wallet-generator --seed "word1 word2 word3 ... word24" --simple

# Create a wallet with custom subwallet ID
./ton-wallet-generator --generate --subwallet 0

# Pipe examples - extract specific fields
./ton-wallet-generator --generate | grep "^Address:"
./ton-wallet-generator --generate | awk -F': ' '/^Address:/ {print $2}'
./ton-wallet-generator --generate --simple | xargs -I {} echo "New wallet: {}"
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
./ton-wallet-generator --generate --simple
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
ADDRESS=$(./ton-wallet-generator --generate --simple 2>&1 | tail -1)

# Get address from existing seed for scripting
./ton-wallet-generator --seed "your 24 words here" --simple

# Generate wallet and filter output for parsing
./ton-wallet-generator --generate | grep "^Address:" | cut -d' ' -f2

# Create multiple wallets with different subwallet IDs from same seed
for i in 0 1 2; do
  ./ton-wallet-generator --seed "your seed" --subwallet $i --simple
done
```

### Creating Multiple Wallets

To create multiple wallets from the same seed phrase, run the program multiple times with:

- Same seed phrase
- Different subwallet IDs (e.g., 0, 1, 2, 3...)
- Same or different wallet versions

### Example

```
TON Wallet Generator
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
This is an offline wallet generator - no network connection required!
```

## Security Notes

- **Never share your seed phrase** - Anyone with access to it can control all wallets derived from it
- **Store your seed phrase securely** - Write it down and keep it in a safe place
- **Private keys are sensitive** - The program displays private keys for convenience, but never share them
- Different subwallet IDs create completely different wallet addresses from the same seed

## Technical Details

- **Completely offline** - No network dependencies
- Uses PBKDF2 with HMAC-SHA512 for key derivation (TON standard)
- Follows TON's standard mnemonic-to-key conversion
- Supports all major wallet contract versions (V3R1, V3R2, V4R1, V4R2, V5R1 Final)
- Network-aware address generation (Mainnet/Testnet)
- Default subwallet ID: 698983191
- Secure key derivation using TON's official method
