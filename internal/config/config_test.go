package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlContent := `
listen_addr: ":8443"
services:
  "service1.local": "localhost:8081"
  "service2.local": "localhost:8082"
tls:
  cert_file: "certs/server.crt"
  key_file: "certs/server.key"
  ca_cert_file: "certs/ca.crt"
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up after test

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Load the config
	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Assertions
	if cfg.ListenAddr != ":8443" {
		t.Errorf("Expected ListenAddr ':8443', got '%s'", cfg.ListenAddr)
	}

	expectedBackends := map[string]string{
		"service1.local": "localhost:8081",
		"service2.local": "localhost:8082",
	}

	if !reflect.DeepEqual(cfg.Services, expectedBackends) {
		t.Errorf("Expected Backends %v, got %v", expectedBackends, cfg.Services)
	}

	if cfg.TLS.CertFile != "certs/server.crt" {
		t.Errorf("Expected TLS.CertFile 'certs/server.crt', got '%s'", cfg.TLS.CertFile)
	}
	if cfg.TLS.KeyFile != "certs/server.key" {
		t.Errorf("Expected TLS.KeyFile 'certs/server.key', got '%s'", cfg.TLS.KeyFile)
	}
	if cfg.TLS.CACertFile != "certs/ca.crt" {
		t.Errorf("Expected TLS.CACertFile 'certs/ca.crt', got '%s'", cfg.TLS.CACertFile)
	}
}
