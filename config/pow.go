//go:build pow

package config

import (
	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/core/types"
)

func newConsensus(_ *Config) types.ConsesusProtocol {
	return pow.New()
}
