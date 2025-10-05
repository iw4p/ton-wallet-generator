package main

import (
	"bufio"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"golang.org/x/crypto/pbkdf2"
)

// WalletConfig holds all configuration for wallet generation
type WalletConfig struct {
	Seed            []string
	NetworkGlobalID int32
	NetworkName     string
	Version         wallet.VersionConfig
	VersionName     string
	SubwalletID     uint32
}

// KeyPair holds the derived cryptographic keys
type KeyPair struct {
	Private ed25519.PrivateKey
	Public  ed25519.PublicKey
}

// WalletInfo contains the complete wallet information
type WalletInfo struct {
	Config  WalletConfig
	Keys    KeyPair
	Address *address.Address
}

// CLIFlags holds command-line arguments
type CLIFlags struct {
	Generate   bool
	Seed       string
	Network    string
	Version    string
	Subwallet  int
	SimpleMode bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flags := parseCLIFlags()

	// Check if running in CLI mode (any flag provided)
	cliMode := flags.Generate || flags.Seed != "" || flags.Network != "" || flags.Version != "" || flags.Subwallet >= 0

	if !cliMode {
		// Interactive mode
		printHeader()
		scanner := bufio.NewScanner(os.Stdin)
		config, err := collectWalletConfig(scanner)
		if err != nil {
			return err
		}

		walletInfo, err := generateWallet(config)
		if err != nil {
			return err
		}

		displayWalletInfo(walletInfo)
		return nil
	}

	// CLI mode - clean output for piping
	config, err := buildConfigFromFlags(flags)
	if err != nil {
		return err
	}

	walletInfo, err := generateWallet(config)
	if err != nil {
		return err
	}

	if flags.SimpleMode {
		fmt.Println(walletInfo.Address.String())
	} else {
		displayWalletInfoCLI(walletInfo)
	}

	return nil
}

func parseCLIFlags() CLIFlags {
	generate := flag.Bool("generate", false, "Generate a new seed phrase")
	seed := flag.String("seed", "", "Seed phrase (24 words, space-separated)")
	network := flag.String("network", "", "Network: mainnet or testnet (default: mainnet)")
	version := flag.String("version", "", "Wallet version: v3r1, v3r2, v4r1, v4r2, v5r1beta, v5r1final (default: v4r2)")
	subwallet := flag.Int("subwallet", -1, "Subwallet ID (default: 698983191 for v3/v4, 0 for v5)")
	simple := flag.Bool("simple", false, "Simple output mode (only show address)")

	flag.Parse()

	return CLIFlags{
		Generate:   *generate,
		Seed:       *seed,
		Network:    *network,
		Version:    *version,
		Subwallet:  *subwallet,
		SimpleMode: *simple,
	}
}

func buildConfigFromFlags(flags CLIFlags) (WalletConfig, error) {
	var seed []string

	// Handle seed
	if flags.Generate {
		seed = wallet.NewSeed()
		// Always output generated seed (even in simple mode) to stderr with clear label
		fmt.Fprintf(os.Stderr, "# Seed phrase (save securely!):\n")
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(seed, " "))
		if !flags.SimpleMode {
			fmt.Fprintln(os.Stderr)
		}
	} else if flags.Seed != "" {
		seed = strings.Fields(flags.Seed)
		if len(seed) != 24 {
			return WalletConfig{}, fmt.Errorf("seed phrase must contain exactly 24 words, got %d", len(seed))
		}
	} else {
		return WalletConfig{}, fmt.Errorf("either --generate or --seed must be provided")
	}

	// Handle network
	var networkID int32
	networkID = wallet.MainnetGlobalID
	networkName := "Mainnet"
	if flags.Network != "" {
		switch strings.ToLower(flags.Network) {
		case "mainnet", "main":
			networkID = wallet.MainnetGlobalID
			networkName = "Mainnet"
		case "testnet", "test":
			networkID = wallet.TestnetGlobalID
			networkName = "Testnet"
		default:
			return WalletConfig{}, fmt.Errorf("invalid network: %s (use mainnet or testnet)", flags.Network)
		}
	}

	// Handle version
	versionStr := "v4r2"
	if flags.Version != "" {
		versionStr = strings.ToLower(flags.Version)
	}

	var walletVersion wallet.VersionConfig
	var versionName string
	var isV5 bool

	switch versionStr {
	case "v3r1":
		walletVersion = wallet.V3R1
		versionName = "V3R1"
	case "v3r2":
		walletVersion = wallet.V3R2
		versionName = "V3R2"
	case "v4r1":
		walletVersion = wallet.V4R1
		versionName = "V4R1"
	case "v4r2":
		walletVersion = wallet.V4R2
		versionName = "V4R2"
	case "v5r1beta", "v5beta":
		config := wallet.ConfigV5R1Beta{NetworkGlobalID: networkID}
		walletVersion = config
		versionName = "V5R1 Beta"
		isV5 = true
	case "v5r1final", "v5r1", "v5":
		config := wallet.ConfigV5R1Final{NetworkGlobalID: networkID}
		walletVersion = config
		versionName = "V5R1 Final"
		isV5 = true
	default:
		return WalletConfig{}, fmt.Errorf("invalid version: %s", flags.Version)
	}

	// Handle subwallet ID
	subwalletID := uint32(wallet.DefaultSubwallet)
	if isV5 {
		subwalletID = 0
	}
	if flags.Subwallet >= 0 {
		subwalletID = uint32(flags.Subwallet)
	}

	return WalletConfig{
		Seed:            seed,
		NetworkGlobalID: networkID,
		NetworkName:     networkName,
		Version:         walletVersion,
		VersionName:     versionName,
		SubwalletID:     subwalletID,
	}, nil
}

// ═══════════════════════════════════════════════════════════════════
// UI & Display Functions
// ═══════════════════════════════════════════════════════════════════

func printHeader() {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║     TON Wallet Manager v1.0        ║")
	fmt.Println("║     Offline Wallet Generator       ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println()
}

func displayWalletInfo(info WalletInfo) {
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║            Wallet Created Successfully!                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  Network:      %s\n", info.Config.NetworkName)
	fmt.Printf("  Version:      %s\n", info.Config.VersionName)
	fmt.Printf("  Subwallet ID: %d\n", info.Config.SubwalletID)
	fmt.Println()
	fmt.Printf("  Address:      %s\n", info.Address.String())
	fmt.Println()
	fmt.Printf("  Public Key:   %x\n", info.Keys.Public)
	fmt.Printf("  Private Key:  %x\n", info.Keys.Private)
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("  ⚠️  SECURITY REMINDERS:")
	fmt.Println("  • Store your seed phrase securely offline")
	fmt.Println("  • Never share your private key or seed phrase")
	fmt.Println("  • This generator works completely offline")
	fmt.Println()
}

func displayWalletInfoCLI(info WalletInfo) {
	fmt.Printf("Address: %s\n", info.Address.String())
	fmt.Printf("Network: %s\n", info.Config.NetworkName)
	fmt.Printf("Version: %s\n", info.Config.VersionName)
	fmt.Printf("Subwallet: %d\n", info.Config.SubwalletID)
	fmt.Printf("PublicKey: %x\n", info.Keys.Public)
	fmt.Printf("PrivateKey: %x\n", info.Keys.Private)
}

// ═══════════════════════════════════════════════════════════════════
// Input Collection Functions
// ═══════════════════════════════════════════════════════════════════

func readMultiLineSeed(scanner *bufio.Scanner) ([]string, error) {
	var allWords []string
	emptyLineCount := 0

	fmt.Println("(Press Enter twice when done, or paste all at once)")
	fmt.Print("> ")

	for {
		if !scanner.Scan() {
			if len(allWords) > 0 {
				break
			}
			return nil, fmt.Errorf("failed to read seed phrase")
		}

		line := scanner.Text()

		// Check if line is empty
		if strings.TrimSpace(line) == "" {
			emptyLineCount++
			// If we have words and user pressed Enter twice, we're done
			if len(allWords) > 0 && emptyLineCount >= 2 {
				break
			}
			// If we already have 24 words and one empty line, we're done
			if len(allWords) >= 24 && emptyLineCount >= 1 {
				break
			}
			continue
		}

		// Reset empty line counter when we get actual input
		emptyLineCount = 0

		// Parse the line: replace common separators with spaces, then split
		line = strings.ReplaceAll(line, ",", " ")
		line = strings.ReplaceAll(line, "\t", " ")
		line = strings.ReplaceAll(line, ";", " ")
		line = strings.ReplaceAll(line, "|", " ")
		line = strings.ReplaceAll(line, ".", " ")

		// Split by whitespace and add non-empty words
		words := strings.Fields(line)
		for _, word := range words {
			word = strings.ToLower(strings.TrimSpace(word))
			if word != "" {
				allWords = append(allWords, word)
			}
		}

		// Show progress
		if len(allWords) > 0 && len(allWords) < 24 {
			fmt.Printf("  (%d/24 words) ", len(allWords))
		}

		// If we have exactly 24 words, we're done
		if len(allWords) >= 24 {
			break
		}

		// Continue reading
		fmt.Print("> ")
	}

	// Take only first 24 words if user pasted more
	if len(allWords) > 24 {
		allWords = allWords[:24]
	}

	return allWords, nil
}

func collectWalletConfig(scanner *bufio.Scanner) (WalletConfig, error) {
	seed, err := promptForSeed(scanner)
	if err != nil {
		return WalletConfig{}, err
	}

	networkID, networkName, err := promptForNetwork(scanner)
	if err != nil {
		return WalletConfig{}, err
	}

	version, versionName, versionChoice, err := promptForVersion(scanner, networkID)
	if err != nil {
		return WalletConfig{}, err
	}

	isV5 := versionChoice == "5" || versionChoice == "6"
	subwalletID, err := promptForSubwallet(scanner, isV5)
	if err != nil {
		return WalletConfig{}, err
	}

	return WalletConfig{
		Seed:            seed,
		NetworkGlobalID: networkID,
		NetworkName:     networkName,
		Version:         version,
		VersionName:     versionName,
		SubwalletID:     subwalletID,
	}, nil
}

func promptForSeed(scanner *bufio.Scanner) ([]string, error) {
	fmt.Print("Generate new seed phrase? (y/n): ")
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read input")
	}

	generateNew := strings.ToLower(strings.TrimSpace(scanner.Text())) == "y"

	if generateNew {
		seed := wallet.NewSeed()
		fmt.Println("\n⚠️  Your new seed phrase (SAVE THIS SECURELY):")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println(strings.Join(seed, " "))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		return seed, nil
	}

	fmt.Println("\nEnter your seed phrase (24 words in any format):")
	fmt.Println("  • Accepts spaces, commas, tabs, or line breaks")
	fmt.Println("  • Paste multiple lines or single line - both work!")
	fmt.Println()

	seed, err := readMultiLineSeed(scanner)
	if err != nil {
		return nil, err
	}

	if len(seed) != 24 {
		return nil, fmt.Errorf("expected 24 words, got %d", len(seed))
	}

	fmt.Printf("✓ Successfully parsed %d words\n\n", len(seed))
	return seed, nil
}

func promptForNetwork(scanner *bufio.Scanner) (int32, string, error) {
	fmt.Println("\nSelect network:")
	fmt.Println("  1. Mainnet")
	fmt.Println("  2. Testnet")
	fmt.Print("\nChoice (1-2): ")

	if !scanner.Scan() {
		return 0, "", fmt.Errorf("failed to read network choice")
	}

	choice := strings.TrimSpace(scanner.Text())
	switch choice {
	case "1":
		fmt.Println("✓ Using Mainnet")
		fmt.Println()
		return wallet.MainnetGlobalID, "Mainnet", nil
	case "2":
		fmt.Println("✓ Using Testnet")
		fmt.Println()
		return wallet.TestnetGlobalID, "Testnet", nil
	default:
		fmt.Println("⚠ Invalid choice, defaulting to Mainnet")
		fmt.Println()
		return wallet.MainnetGlobalID, "Mainnet", nil
	}
}

func promptForVersion(scanner *bufio.Scanner, networkID int32) (wallet.VersionConfig, string, string, error) {
	fmt.Println("Available wallet versions:")
	fmt.Println("  1. V3R1")
	fmt.Println("  2. V3R2")
	fmt.Println("  3. V4R1")
	fmt.Println("  4. V4R2 (recommended)")
	fmt.Println("  5. V5R1 Beta")
	fmt.Println("  6. V5R1 Final")
	fmt.Print("\nChoice (1-6): ")

	if !scanner.Scan() {
		return nil, "", "", fmt.Errorf("failed to read version choice")
	}

	choice := strings.TrimSpace(scanner.Text())
	fmt.Println()

	switch choice {
	case "1":
		return wallet.V3R1, "V3R1", choice, nil
	case "2":
		return wallet.V3R2, "V3R2", choice, nil
	case "3":
		return wallet.V4R1, "V4R1", choice, nil
	case "4":
		return wallet.V4R2, "V4R2", choice, nil
	case "5":
		config := wallet.ConfigV5R1Beta{NetworkGlobalID: networkID}
		return config, "V5R1 Beta", choice, nil
	case "6":
		config := wallet.ConfigV5R1Final{NetworkGlobalID: networkID}
		return config, "V5R1 Final", choice, nil
	default:
		fmt.Println("⚠ Invalid choice, defaulting to V4R2")
		fmt.Println()
		return wallet.V4R2, "V4R2", "4", nil
	}
}

func promptForSubwallet(scanner *bufio.Scanner, isV5 bool) (uint32, error) {
	defaultSubwallet := uint32(wallet.DefaultSubwallet)
	defaultDesc := "698983191 (standard)"

	if isV5 {
		defaultSubwallet = 0
		defaultDesc = "0 (W5 compatible)"
	}

	fmt.Println("Subwallet ID options:")
	fmt.Printf("  • Press Enter for default (%s)\n", defaultDesc)
	fmt.Println("  • Enter custom number for multiple wallets from same seed")
	fmt.Print("\nSubwallet ID: ")

	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to read subwallet ID")
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		fmt.Printf("✓ Using default subwallet ID (%s)\n", defaultDesc)
		return defaultSubwallet, nil
	}

	parsed, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid subwallet ID: %v", err)
	}

	subwalletID := uint32(parsed)
	fmt.Printf("✓ Using custom subwallet ID (%d)\n", subwalletID)
	return subwalletID, nil
}

// ═══════════════════════════════════════════════════════════════════
// Wallet Generation Functions
// ═══════════════════════════════════════════════════════════════════

func generateWallet(config WalletConfig) (WalletInfo, error) {
	keys, err := deriveKeysFromSeed(config.Seed)
	if err != nil {
		return WalletInfo{}, fmt.Errorf("failed to derive keys: %w", err)
	}

	addr, err := wallet.AddressFromPubKey(keys.Public, config.Version, config.SubwalletID)
	if err != nil {
		return WalletInfo{}, fmt.Errorf("failed to create wallet address: %w", err)
	}

	return WalletInfo{
		Config:  config,
		Keys:    keys,
		Address: addr,
	}, nil
}

func deriveKeysFromSeed(seed []string) (KeyPair, error) {
	seedPhrase := strings.Join(seed, " ")

	// Create HMAC-SHA512 hash of the seed phrase
	mac := hmac.New(sha512.New, []byte(seedPhrase))
	mac.Write([]byte(""))
	hash := mac.Sum(nil)

	// Use PBKDF2 with the hash as input (TON standard derivation)
	derivedKey := pbkdf2.Key(hash, []byte("TON default seed"), 100000, 32, sha512.New)

	// Generate ed25519 key pair from derived key
	privateKey := ed25519.NewKeyFromSeed(derivedKey)

	// Extract public key bytes directly from private key
	publicKeyBytes := privateKey[32:]
	publicKey := ed25519.PublicKey(publicKeyBytes)

	return KeyPair{
		Private: privateKey,
		Public:  publicKey,
	}, nil
}
