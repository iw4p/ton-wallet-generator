package ui

import (
	"fmt"

	"ton-wallet-manager/internal/wallet"
)

// PrintHeader displays the application header
func PrintHeader() {
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║     TON Wallet Manager v1.0       ║")
	fmt.Println("║     Offline Wallet Management      ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println()
}

// DisplayWalletInfo displays wallet information in interactive mode
func DisplayWalletInfo(info wallet.Info) {
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

// DisplayWalletInfoCLI displays wallet information in CLI mode
func DisplayWalletInfoCLI(info wallet.Info) {
	fmt.Printf("Address: %s\n", info.Address.String())
	fmt.Printf("Network: %s\n", info.Config.NetworkName)
	fmt.Printf("Version: %s\n", info.Config.VersionName)
	fmt.Printf("Subwallet: %d\n", info.Config.SubwalletID)
	fmt.Printf("PublicKey: %x\n", info.Keys.Public)
	fmt.Printf("PrivateKey: %x\n", info.Keys.Private)
}

// DisplaySimple displays only the wallet address (for simple mode)
func DisplaySimple(info wallet.Info) {
	fmt.Println(info.Address.String())
}
