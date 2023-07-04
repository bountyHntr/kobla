package main

import (
	"log"
	"path2perpetuity/blockchain/core/chain"
)

func main() {
	bc, err := chain.New()
	if err != nil {
		log.Fatalf("create blockchian: %s", err)
	}

	bc.AddBlock([]byte("Send 1 BTC to Ivan"))
	bc.AddBlock([]byte("Send 2 more BTC to Ivan"))
}
