package cli

import "flag"

// Flags holds command-line arguments
type Flags struct {
	Generate   bool
	Seed       string
	Network    string
	Version    string
	Subwallet  int
	SimpleMode bool
}

// ParseFlags parses command-line flags and returns a Flags struct
func ParseFlags() Flags {
	generate := flag.Bool("generate", false, "Generate a new seed phrase")
	seed := flag.String("seed", "", "Seed phrase (24 words, space-separated)")
	network := flag.String("network", "", "Network: mainnet or testnet (default: mainnet)")
	version := flag.String("version", "", "Wallet version: v3r1, v3r2, v4r1, v4r2, v5r1beta, v5r1final (default: v4r2)")
	subwallet := flag.Int("subwallet", -1, "Subwallet ID (default: 698983191 for v3/v4, 0 for v5)")
	simple := flag.Bool("simple", false, "Simple output mode (only show address)")

	flag.Parse()

	return Flags{
		Generate:   *generate,
		Seed:       *seed,
		Network:    *network,
		Version:    *version,
		Subwallet:  *subwallet,
		SimpleMode: *simple,
	}
}

// IsCLIMode returns true if any CLI flags were provided
func (f Flags) IsCLIMode() bool {
	return f.Generate || f.Seed != "" || f.Network != "" || f.Version != "" || f.Subwallet >= 0
}
