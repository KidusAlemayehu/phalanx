package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/KidusAlemayehu/phalanx/internal/config"
	"github.com/KidusAlemayehu/phalanx/internal/proxy"
)

var exitFunc = os.Exit

func main() {
	defaultConfigFilePath := config.DEFAULT_CONFIG_FILE_PATH
	configFilePath := flag.String("conf", defaultConfigFilePath, "Path to proxy config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFilePath)

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
