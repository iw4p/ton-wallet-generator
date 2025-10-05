package cli

import (
	"fmt"
	"os"
	"strings"

	"ton-wallet-manager/internal/wallet"

	tonwallet "github.com/xssnick/tonutils-go/ton/wallet"
)

// BuildConfig builds a wallet.Config from CLI flags
func BuildConfig(flags Flags) (wallet.Config, error) {
	var seed []string

	// Handle seed
	if flags.Generate {
		seed = tonwallet.NewSeed()
		// Always output generated seed (even in simple mode) to stderr with clear label
		fmt.Fprintf(os.Stderr, "# Seed phrase (save securely!):\n")
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(seed, " "))
		if !flags.SimpleMode {
			fmt.Fprintln(os.Stderr)
		}
	} else if flags.Seed != "" {
		seed = strings.Fields(flags.Seed)
		if len(seed) != 24 {
			return wallet.Config{}, fmt.Errorf("seed phrase must contain exactly 24 words, got %d", len(seed))
		}
	} else {
		return wallet.Config{}, fmt.Errorf("either --generate or --seed must be provided")
	}

	// Handle network
	var networkID int32
	networkID = tonwallet.MainnetGlobalID
	networkName := "Mainnet"
	if flags.Network != "" {
		switch strings.ToLower(flags.Network) {
		case "mainnet", "main":
			networkID = tonwallet.MainnetGlobalID
			networkName = "Mainnet"
		case "testnet", "test":
			networkID = tonwallet.TestnetGlobalID
			networkName = "Testnet"
		default:
			return wallet.Config{}, fmt.Errorf("invalid network: %s (use mainnet or testnet)", flags.Network)
		}
	}

	// Handle version
	versionStr := "v4r2"
	if flags.Version != "" {
		versionStr = strings.ToLower(flags.Version)
	}

	var walletVersion tonwallet.VersionConfig
	var versionName string
	var isV5 bool

	switch versionStr {
	case "v3r1":
		walletVersion = tonwallet.V3R1
		versionName = "V3R1"
	case "v3r2":
		walletVersion = tonwallet.V3R2
		versionName = "V3R2"
	case "v4r1":
		walletVersion = tonwallet.V4R1
		versionName = "V4R1"
	case "v4r2":
		walletVersion = tonwallet.V4R2
		versionName = "V4R2"
	case "v5r1beta", "v5beta":
		config := tonwallet.ConfigV5R1Beta{NetworkGlobalID: networkID}
		walletVersion = config
		versionName = "V5R1 Beta"
		isV5 = true
	case "v5r1final", "v5r1", "v5":
		config := tonwallet.ConfigV5R1Final{NetworkGlobalID: networkID}
		walletVersion = config
		versionName = "V5R1 Final"
		isV5 = true
	default:
		return wallet.Config{}, fmt.Errorf("invalid version: %s", flags.Version)
	}

	// Handle subwallet ID
	subwalletID := uint32(tonwallet.DefaultSubwallet)
	if isV5 {
		subwalletID = 0
	}
	if flags.Subwallet >= 0 {
		subwalletID = uint32(flags.Subwallet)
	}

	return wallet.Config{
		Seed:            seed,
		NetworkGlobalID: networkID,
		NetworkName:     networkName,
		Version:         walletVersion,
		VersionName:     versionName,
		SubwalletID:     subwalletID,
	}, nil
}
