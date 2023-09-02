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
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

type Validator struct {
	Url     string `yaml:"URL"`
	Address string `yaml:"Address"`
}

type Config struct {
	Url          string      `yaml:"URL"`
	ListeningUrl string      `yaml:"ListeningURL"`
	Genesis      bool        `yaml:"Genesis"`
	DbPath       string      `yaml:"DBPath"`
	Consensus    string      `yaml:"Consensus"`
	Nodes        []string    `yaml:"Nodes"`
	Validators   []Validator `yaml:"Validators"`
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

	log.Infof("path to database: %s", cfg.DbPath)
	log.Infof("consensus protocol: %s", cfg.Consensus)
	log.Infof("node URL: %s", cfg.Url)
	log.Infof("node listening URL: %s", cfg.ListeningUrl)

	if cfg.Genesis {
		log.Infof("launching new blockchain")
	}

	return chain.Config{
		DBPath:       cfg.DbPath,
		Consensus:    parseConsensusProtocol(&cfg),
		Url:          cfg.Url,
		ListeningUrl: cfg.ListeningUrl,
		Nodes:        parseNodes(&cfg),
		Genesis:      cfg.Genesis,
	}
}

func parseNodes(cfg *Config) []string {
	nodes := make([]string, 0, len(cfg.Nodes))

	for _, node := range cfg.Nodes {
		if node == cfg.Url {
			continue
		}

		nodes = append(nodes, node)
		log.Infof("node: %s", node)
	}

	return nodes
}

func parseConsensusProtocol(cfg *Config) (consensus types.ConsesusProtocol) {
	switch cfg.Consensus {
	case "", "PoW":
		consensus = pow.New()
	case "PoA":
		log.Infof("%d validators:", len(cfg.Validators))

		validators := make([]poa.Validator, 0, len(cfg.Nodes))
		for _, v := range cfg.Validators {
			log.Infof("validator url: %s", v.Url)
			log.Infof("validator address: %s", v.Address)

			validators = append(validators, poa.Validator{
				Url:     v.Url,
				Address: types.AddressFromString(v.Address),
			})
		}

		consensus = poa.New(validators)
	}

	return
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
