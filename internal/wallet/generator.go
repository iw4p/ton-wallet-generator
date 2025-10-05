package wallet

import (
	"fmt"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

// Generate creates a wallet from the given configuration
func Generate(config Config) (Info, error) {
	keys, err := DeriveKeysFromSeed(config.Seed)
	if err != nil {
		return Info{}, fmt.Errorf("failed to derive keys: %w", err)
	}

	addr, err := wallet.AddressFromPubKey(keys.Public, config.Version, config.SubwalletID)
	if err != nil {
		return Info{}, fmt.Errorf("failed to create wallet address: %w", err)
	}

	return Info{
		Config:  config,
		Keys:    keys,
		Address: addr,
	}, nil
}
