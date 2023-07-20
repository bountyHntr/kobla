package main

import (
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/tui"
	"log"
)

const dbPath = "./.data"

func main() {
	cfg := chain.Config{
		DBPath:    dbPath,
		Consensus: pow.New(),
	}

	bc, err := chain.New(&cfg)
	if err != nil {
		log.Fatalf("new blockchain: %s", err)
	}

	if err := tui.Run(bc); err != nil {
		log.Fatalf("app: %s", err)
	}
}
