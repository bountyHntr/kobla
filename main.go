package main

import (
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/tui"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
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

	go func() {
		for i := 0; i < 100; i++ {
			if err := bc.AddBlock([]byte(strconv.Itoa(i))); err != nil {
				log.Fatalf("add block: %s", err)
			}

			time.Sleep(time.Second)
		}
	}()

	log.Info("run")

	if err := tui.Run(bc); err != nil {
		log.Fatalf("app: %s", err)
	}
}
