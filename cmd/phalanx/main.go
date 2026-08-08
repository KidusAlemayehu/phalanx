package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KidusAlemayehu/phalanx/internal/config"
	"github.com/KidusAlemayehu/phalanx/internal/proxy"
)

var exitFunc = os.Exit

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		exitFunc(1)
	}

	log.Printf("Proxy server started at : %s", cfg.ListenAddr)

	err = proxy.Start(cfg)

	if err != nil {
		log.Fatalf("Proxy Server crashed : %v", err)
	}
}
