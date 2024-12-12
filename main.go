package main

import (
	"log"
	"trustify/config"
	"trustify/network"
)

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}

	node := network.NewNode(cfg)

	node.Start()
}
