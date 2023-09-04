package config

import (
	"flag"
	"kobla/blockchain/core/chain"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Url        string   `yaml:"URL"`
	Genesis    bool     `yaml:"Genesis"`
	DbPath     string   `yaml:"DBPath"`
	Nodes      []string `yaml:"Nodes"`
	Validators []string `yaml:"Validators"`
	PrivateKey string   `yaml:"PrivateKey"`
}

func Build() chain.Config {
	cfgPath := flag.String("config", "config.yaml", "path to yaml config file")
	flag.Parse()

	var cfg Config
	if err := parseConfigFromFile(*cfgPath, &cfg); err != nil {
		log.Fatalf("failed to read config file %s: %s", *cfgPath, err)
	}

	log.Infof("path to database: %s", cfg.DbPath)
	log.Infof("node URL: %s", cfg.Url)

	if cfg.Genesis {
		log.Infof("launching new blockchain")
	}

	return chain.Config{
		DBPath:    cfg.DbPath,
		Consensus: newConsensus(&cfg),
		Url:       cfg.Url,
		Nodes:     cfg.Nodes,
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
