package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BridgeAddress string `yaml:"bridgeAddress"`
	ClientID      string `yaml:"clientID"`
}

// Read the config from the file system
func Read() (*Config, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("unable to find the configuration directory: %w", err)
	}

	cfg := filepath.Join(cfgDir, "huectl", "config.yml")
	if _, err = os.Stat(cfg); os.IsNotExist(err) {
		log.Warnf("Expected file at %v", cfg)
		return nil, err
	}

	f, err := os.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	var config Config
	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}
	log.Debugf("Loaded config: %+v", config)

	return &config, nil
}