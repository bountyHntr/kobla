package main

import (
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/consensus/pow"
	"log"
)

func main() {
	cfg := chain.Config{
		DBPath:    "./.data",
		Consensus: pow.New(),
	}

	bc, err := chain.New(&cfg)
	if err != nil {
		log.Fatalf("create blockchian: %s", err)
	}

	block, err := bc.BlockByNumber(2)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v", block)

	// bc.AddBlock([]byte("Send 1 BTC to Ivan"))
	// bc.AddBlock([]byte("Send 2 more BTC to Ivan"))
}
