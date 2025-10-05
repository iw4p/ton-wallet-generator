package wallet

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// DeriveKeysFromSeed derives an ed25519 key pair from a seed phrase using TON standard derivation
func DeriveKeysFromSeed(seed []string) (KeyPair, error) {
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
