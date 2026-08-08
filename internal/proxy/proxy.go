package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/KidusAlemayehu/phalanx/internal/cert"
	"github.com/KidusAlemayehu/phalanx/internal/config"
)

func Start(cfg *config.Config) error {
	tlsConfig, err := cert.LoadTLSConfig(cfg)

	if err != nil {
		return err
	}

	listener, err := tls.Listen("tcp", cfg.ListenAddr, tlsConfig)

	if err != nil {
		return fmt.Errorf("Couldn't listen to tcp addr")
	}

	defer listener.Close()

	for {
		clientConn, err := listener.Accept()

		if err != nil {
			log.Printf("Error accepting tcp connections: %v", err)
			continue
		}
		go handleTCPConnection(clientConn, cfg)
	}
}

func handleTCPConnection(clientConn net.Conn, cfg *config.Config) {
	defer clientConn.Close()

	backendConn, err := net.Dial("tcp", cfg.BackendAddr)

	if err != nil {
		log.Fatalf("Failed to reach backend app at %s : %v", cfg.BackendAddr, err)
		return
	}

	defer backendConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, err := io.Copy(backendConn, clientConn)
		if err != nil {
			log.Printf("stream connection to backend error %v", err)
		}

		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(clientConn, backendConn)

		if err != nil {
			log.Printf("stream connection to client error %v", err)
		}

		done <- struct{}{}
	}()

}
