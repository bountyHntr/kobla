package main

import (
	"flag"
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/tui"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	cfg := buildConfig()

	bc, err := chain.New(&cfg)
	if err != nil {
		log.Fatalf("new blockchain: %s", err)
	}

	if err := tui.Run(bc); err != nil {
		log.Fatalf("app: %s", err)
	}
}

func buildConfig() chain.Config {

	url := flag.String("url", "localhost:8090", "address where the node is listening")
	genesis := flag.Bool("genesis", false, "flag that indicates that a new network is being created")
	syncNode := flag.String("sync_node", "", "flag that indicates that a new network is being created")
	dbPath := flag.String("db_path", "./.data", "node database path")

	flag.Parse()

	return chain.Config{
		DBPath:    *dbPath,
		Consensus: pow.New(),
		URL:       *url,
		SyncNode:  *syncNode,
		Genesis:   *genesis,
	}
}
