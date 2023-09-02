package main

import (
	"flag"
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/consensus/poa"
	"kobla/blockchain/core/consensus/pow"
	"kobla/blockchain/core/types"
	"kobla/blockchain/tui"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Node struct {
	Url     string `yaml:"URL"`
	Address string `yaml:"Address"`
}

type Config struct {
	Url       string `yaml:"URL"`
	Genesis   bool   `yaml:"Genesis"`
	DbPath    string `yaml:"DBPath"`
	Consensus string `yaml:"Consensus"`
	Nodes     []Node `yaml:"Nodes"`
}

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

	cfgPath := flag.String("cfg", "config.yaml", "path to yaml config file")
	flag.Parse()

	var cfg Config
	if err := parseConfigFromFile(*cfgPath, &cfg); err != nil {
		log.Fatalf("failed to read config file %s: %s", *cfgPath, err)
	}

	nodes := make([]string, len(cfg.Nodes))
	for _, node := range cfg.Nodes {
		nodes = append(nodes, node.Url)
	}

	var consensus types.ConsesusProtocol
	switch cfg.Consensus {
	case "pow", "":
		consensus = pow.New()
	case "poa":
		validators := make([]poa.Validator, len(cfg.Nodes))
		for _, node := range cfg.Nodes {
			validators = append(validators, poa.Validator{
				Url:     node.Url,
				Address: types.AddressFromString(node.Address),
			})
		}

		consensus = poa.New(validators)
	}

	return chain.Config{
		DBPath:    cfg.DbPath,
		Consensus: consensus,
		URL:       cfg.Url,
		Nodes:     nodes,
		Genesis:   cfg.Genesis,
	}
}

func parseConfigFromFile(fileName string, cfg interface{}) error {
	rawCfg, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return parseConfigRaw(rawCfg, cfg)
}

func parseConfigRaw(rawCfg []byte, cfg interface{}) error {
	err := yaml.Unmarshal(rawCfg, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config file")
	}
	return nil
}
