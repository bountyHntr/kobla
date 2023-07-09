package main

import (
	"kobla/blockchain/tui"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := tui.Run(); err != nil {
		log.Fatalf("app: %s", err)
	}
}
