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

	tlsConn, ok := clientConn.(*tls.Conn)
	if !ok {
		log.Printf("Expected a TLS connection")
		return
	}

	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}

	serverName := tlsConn.ConnectionState().ServerName
	serviceAddr, ok := cfg.Services[serverName]
	if !ok {
		log.Printf("No backend found for SNI: %s", serverName)
		return
	}

	backendConn, err := net.Dial("tcp", serviceAddr)

	if err != nil {
		log.Printf("Failed to reach backend app for SNI %s at %s: %v", serverName, serviceAddr, err)
		return
	}

	defer backendConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, err := io.Copy(backendConn, tlsConn)
		if err != nil {
			log.Printf("stream connection to backend error %v", err)
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(tlsConn, backendConn)

		if err != nil {
			log.Printf("stream connection to client error %v", err)
		}
		done <- struct{}{}
	}()

	<-done
}
