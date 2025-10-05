package ui

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"ton-wallet-manager/internal/wallet"

	tonwallet "github.com/xssnick/tonutils-go/ton/wallet"
)

// CollectWalletConfig collects wallet configuration interactively
func CollectWalletConfig(scanner *bufio.Scanner) (wallet.Config, error) {
	seed, err := promptForSeed(scanner)
	if err != nil {
		return wallet.Config{}, err
	}

	networkID, networkName, err := promptForNetwork(scanner)
	if err != nil {
		return wallet.Config{}, err
	}

	version, versionName, versionChoice, err := promptForVersion(scanner, networkID)
	if err != nil {
		return wallet.Config{}, err
	}

	isV5 := versionChoice == "5" || versionChoice == "6"
	subwalletID, err := promptForSubwallet(scanner, isV5)
	if err != nil {
		return wallet.Config{}, err
	}

	return wallet.Config{
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
		seed := tonwallet.NewSeed()
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
		return tonwallet.MainnetGlobalID, "Mainnet", nil
	case "2":
		fmt.Println("✓ Using Testnet")
		fmt.Println()
		return tonwallet.TestnetGlobalID, "Testnet", nil
	default:
		fmt.Println("⚠ Invalid choice, defaulting to Mainnet")
		fmt.Println()
		return tonwallet.MainnetGlobalID, "Mainnet", nil
	}
}

func promptForVersion(scanner *bufio.Scanner, networkID int32) (tonwallet.VersionConfig, string, string, error) {
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
		return tonwallet.V3R1, "V3R1", choice, nil
	case "2":
		return tonwallet.V3R2, "V3R2", choice, nil
	case "3":
		return tonwallet.V4R1, "V4R1", choice, nil
	case "4":
		return tonwallet.V4R2, "V4R2", choice, nil
	case "5":
		config := tonwallet.ConfigV5R1Beta{NetworkGlobalID: networkID}
		return config, "V5R1 Beta", choice, nil
	case "6":
		config := tonwallet.ConfigV5R1Final{NetworkGlobalID: networkID}
		return config, "V5R1 Final", choice, nil
	default:
		fmt.Println("⚠ Invalid choice, defaulting to V4R2")
		fmt.Println()
		return tonwallet.V4R2, "V4R2", "4", nil
	}
}

func promptForSubwallet(scanner *bufio.Scanner, isV5 bool) (uint32, error) {
	defaultSubwallet := uint32(tonwallet.DefaultSubwallet)
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
