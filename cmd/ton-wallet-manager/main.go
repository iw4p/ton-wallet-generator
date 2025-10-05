package main

import (
	"bufio"
	"fmt"
	"os"

	"ton-wallet-manager/internal/cli"
	"ton-wallet-manager/internal/ui"
	"ton-wallet-manager/internal/wallet"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flags := cli.ParseFlags()

	if !flags.IsCLIMode() {
		// Interactive mode
		return runInteractive()
	}

	// CLI mode - clean output for piping
	return runCLI(flags)
}

func runInteractive() error {
	ui.PrintHeader()

	scanner := bufio.NewScanner(os.Stdin)
	config, err := ui.CollectWalletConfig(scanner)
	if err != nil {
		return err
	}

	walletInfo, err := wallet.Generate(config)
	if err != nil {
		return err
	}

	ui.DisplayWalletInfo(walletInfo)
	return nil
}

func runCLI(flags cli.Flags) error {
	config, err := cli.BuildConfig(flags)
	if err != nil {
		return err
	}

	walletInfo, err := wallet.Generate(config)
	if err != nil {
		return err
	}

	if flags.SimpleMode {
		ui.DisplaySimple(walletInfo)
	} else {
		ui.DisplayWalletInfoCLI(walletInfo)
	}

	return nil
}
