package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configFile = "config.yaml"
)

type Config struct {
	ListenAddr  string    `yaml:"listen_addr"`
	BackendAddr string    `yaml:"backend_addr"`
	TLS         TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	CACertFile string `yaml:"ca_cert_file"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)

	if err != nil {
		return nil, fmt.Errorf("Couldn't read config yaml file")
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, fmt.Errorf("Unable to marshal yaml file data")
	}

	return &cfg, nil
}
