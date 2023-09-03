//go:build poa

package config

import (
	"kobla/blockchain/core/consensus/poa"
	"kobla/blockchain/core/types"

	log "github.com/sirupsen/logrus"
)

func newConsensus(cfg *Config) types.ConsesusProtocol {
	log.Infof("%d validators:", len(cfg.Validators))
	consensus, err := poa.New(cfg.Validators, cfg.PrivateKey)
	if err != nil {
		log.Fatalf("failed to init Proof-of-Authority: %s", err)
	}

	return consensus
}
