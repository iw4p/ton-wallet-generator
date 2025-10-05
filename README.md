# TON Wallet Manager

A simple **offline** Go tool to create multiple TON wallets from a single seed phrase using different wallet versions and subwallet IDs.

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
go build -o ton-wallet-manager
```

## Usage

Run the program:

```bash
./ton-wallet-manager
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

### Creating Multiple Wallets

To create multiple wallets from the same seed phrase, run the program multiple times with:

- Same seed phrase
- Different subwallet IDs (e.g., 0, 1, 2, 3...)
- Same or different wallet versions

### Example

```
TON Wallet Manager
==================

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
