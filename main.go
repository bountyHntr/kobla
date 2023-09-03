package main

import (
	"kobla/blockchain/core/chain"
	"kobla/blockchain/tui"

	"kobla/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	cfg := config.Build()

	bc, err := chain.New(&cfg)
	if err != nil {
		log.Fatalf("new blockchain: %s", err)
	}

	if err := tui.Run(bc); err != nil {
		log.Fatalf("app: %s", err)
	}
}
