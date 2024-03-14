package main

import (
	"flag"
	"github.com/s02190058/warehouse/internal/config"
	"log"
)

var configPath = flag.String("config", "./configs/local.yml", "path to config file")

func main() {
	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}

	_ = cfg
}
