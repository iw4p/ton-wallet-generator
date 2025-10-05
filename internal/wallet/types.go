package wallet

import (
	"crypto/ed25519"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

// Config holds all configuration for wallet generation
type Config struct {
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

// Info contains the complete wallet information
type Info struct {
	Config  Config
	Keys    KeyPair
	Address *address.Address
}
