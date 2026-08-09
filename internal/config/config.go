package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenAddr string    `yaml:"listen_addr"`
	TargetAddr string    `yaml:"target_addr"`
	TLS        TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	CACertFile string `yaml:"ca_cert_file"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	configFile, err := os.Open(configFilePath)

	if err != nil {
		return nil, fmt.Errorf("Couldn't read config yaml file")
	}

	var cfg Config

	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("Critical Error: Failed parsing config file")
	}

	return &cfg, nil
}
