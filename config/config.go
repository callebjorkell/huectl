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
	cfg, err := getConfigName()
	if err != nil {
		return nil, err
	}
	if _, err = os.Stat(cfg); os.IsNotExist(err) {
		return nil, fmt.Errorf("expected file at %v", cfg)
	}

	f, err := os.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}
	defer f.Close()

	var config Config
	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}
	log.Debugf("Loaded config: %+v", config)

	return &config, nil
}

func getConfigName() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("unable to find the configuration directory: %w", err)
	}

	cfg := filepath.Join(cfgDir, "huectl", "config.yml")

	return cfg, nil
}

func (c Config) Write() error {
	cfg, err := getConfigName()
	if err != nil {
		return err
	}

	dir := filepath.Dir(cfg)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return fmt.Errorf("could not create config directory: %w", err)
		}
	}

	f, err := os.Create(cfg)
	if err != nil {
		return fmt.Errorf("unable to open config file for writing: %w", err)
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(c)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	return nil
}
