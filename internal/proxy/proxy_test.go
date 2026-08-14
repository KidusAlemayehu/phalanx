package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/KidusAlemayehu/phalanx/internal/config"
)

// generateCerts creates a self-signed CA, a server cert, and a client cert for testing.
func generateCerts(t *testing.T, dir string) {
	// Simple bash script to generate certs using openssl
	script := `
#!/bin/bash
cd ` + dir + `

# CA
openssl req -x509 -nodes -newkey rsa:2048 -days 1 -keyout ca.key -out ca.crt -subj "/CN=Test CA"

# Server
openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr -subj "/CN=server"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 1 -extfile <(printf "subjectAltName=DNS:service1.local,DNS:service2.local")

# Client
openssl req -newkey rsa:2048 -nodes -keyout client.key -out client.csr -subj "/CN=client"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 1
`
	scriptPath := filepath.Join(dir, "gen.sh")
	err := os.WriteFile(scriptPath, []byte(script), 0755)
	if err != nil {
		t.Fatalf("Failed to write cert gen script: %v", err)
	}
	cmd := exec.Command("bash", scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to generate certs: %v\nOutput: %s", err, out)
	}
}

// dummyBackend starts a simple TCP server that echoes back its given name.
func dummyBackend(t *testing.T, addr, name string) net.Listener {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				// Read until EOF or something
				buf := make([]byte, 1024)
				n, err := c.Read(buf)
				if err != nil && err != io.EOF {
					return
				}
				if n > 0 {
					c.Write([]byte(name))
				}
			}(conn)
		}
	}()
	return l
}

func TestProxyRouting(t *testing.T) {
	// 1. Setup Certs
	certDir := t.TempDir()
	generateCerts(t, certDir)

	caCertPath := filepath.Join(certDir, "ca.crt")
	serverCertPath := filepath.Join(certDir, "server.crt")
	serverKeyPath := filepath.Join(certDir, "server.key")
	clientCertPath := filepath.Join(certDir, "client.crt")
	clientKeyPath := filepath.Join(certDir, "client.key")

	// 2. Setup Backends
	backend1Addr := "127.0.0.1:0" // Let OS pick port
	l1 := dummyBackend(t, backend1Addr, "backend1")
	defer l1.Close()
	backend1Addr = l1.Addr().String()

	backend2Addr := "127.0.0.1:0"
	l2 := dummyBackend(t, backend2Addr, "backend2")
	defer l2.Close()
	backend2Addr = l2.Addr().String()

	// 3. Setup Proxy Config
	proxyAddr := "127.0.0.1:0" // Pick free port
	cfg := &config.Config{
		ListenAddr: proxyAddr,
		Backends: map[string]string{
			"service1.local": backend1Addr,
			"service2.local": backend2Addr,
		},
		TLS: config.TLSConfig{
			CertFile:   serverCertPath,
			KeyFile:    serverKeyPath,
			CACertFile: caCertPath,
		},
	}

	// Re-assign a static port for proxy to easily connect, since Start doesn't return the listener
	cfg.ListenAddr = "127.0.0.1:28443"

	// 4. Start Proxy (non-blocking)
	go func() {
		err := Start(cfg)
		if err != nil {
			fmt.Printf("Proxy start error: %v\n", err)
		}
	}()

	// Wait a moment for proxy to start
	time.Sleep(500 * time.Millisecond)

	// 5. Setup Client TLS Config
	caCert, _ := os.ReadFile(caCertPath)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, _ := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)

	runTest := func(sni, expectedResponse string) {
		clientTLSConfig := &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
			ServerName:   sni,
		}

		conn, err := tls.Dial("tcp", cfg.ListenAddr, clientTLSConfig)
		if err != nil {
			t.Fatalf("Failed to dial proxy for %s: %v", sni, err)
		}
		defer conn.Close()

		_, err = conn.Write([]byte("ping\n"))
		if err != nil {
			t.Fatalf("Failed to write to proxy: %v", err)
		}

		// Wait briefly to allow backend to respond
		time.Sleep(100 * time.Millisecond)

		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			t.Logf("Failed to read from proxy (but continuing to check buffer): %v", err)
		}

		response := string(buf[:n])
		if !strings.Contains(response, expectedResponse) {
			t.Errorf("For SNI %s, expected response containing %s, got '%s'", sni, expectedResponse, response)
		}
	}

	// 6. Run Tests
	runTest("service1.local", "backend1")
	runTest("service2.local", "backend2")
}
