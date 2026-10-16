package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/KidusAlemayehu/phalanx/internal/cert"
	"github.com/KidusAlemayehu/phalanx/internal/config"
)

const (
	defaultHandshaketimeout = 10 * time.Second
	defaultDialTimeout      = 5 * time.Second
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

	if err := clientConn.SetDeadline(time.Now().Add(defaultHandshaketimeout)); err != nil {
		log.Printf("Failed to set deadline for TLS handshake: %v", err)
		return
	}

	if err := tlsConn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}

	if err := clientConn.SetDeadline(time.Time{}); err != nil {
		log.Printf("Failed to reset deadline after TLS handshake: %v", err)
		return
	}

	serverName := tlsConn.ConnectionState().ServerName
	serviceAddr, ok := cfg.Services[serverName]
	if !ok {
		log.Printf("No backend found for SNI: %s", serverName)
		return
	}

	backendConn, err := net.DialTimeout("tcp", serviceAddr, defaultDialTimeout)

	if err != nil {
		log.Printf("Failed to reach backend app for SNI %s at %s: %v", serverName, serviceAddr, err)
		return
	}

	defer backendConn.Close()

	done := make(chan struct{}, 2)

	go func() {
		_, err := io.Copy(backendConn, tlsConn)
		if err != nil {
			if !isExpectedCloseError(err) {
				log.Printf("stream connection to backend error %v", err)
			}
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

func isExpectedCloseError(err error) bool {
	if err == nil {
		return false
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	if err == io.EOF {
		return true
	}

	return false
}
