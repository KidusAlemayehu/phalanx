package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/KidusAlemayehu/phalanx/internal/config"
)

func LoadTLSConfig(cfg *config.Config) (*tls.Config, error) {
	serverCert, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)

	if err != nil {
		return nil, fmt.Errorf("Failed to load server cert file")
	}

	caCert, err := os.ReadFile(cfg.TLS.CACertFile)

	if err != nil {
		return nil, fmt.Errorf("Load to fail ca cert file")
	}

	caCertPool := x509.NewCertPool()

	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("Failed to load CA Cert pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}
	return tlsConfig, nil
}
