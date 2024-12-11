package main

import (
	"log"
	"trustify/config"
	"trustify/network"
)

func main() {
	// Load the configuration file (e.g., config.yaml).
	// Validate the configuration file to ensure all required fields are present.
	// Create a new node by passing the config file as a parameter
	// Handle errors gracefully if node initialization fails.
	// Call the Start method on the node to begin operations like networking, transaction processing, and mining.
	// Maintain an infinite loop to keep the program alive, allowing the node to operate continuously.
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}

	// // Proceed with initializing the node using cfg
	node := network.NewNode(cfg)

	node.Start()
	// hash := "0fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	// enc, _ := hex.DecodeString(hash)
	// print(hex.EncodeToString(enc))

	// // Step 4: Set up graceful shutdown handling.
	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// fmt.Println("Node is running. Press Ctrl+C to exit.")

	// // Step 5: Block until a termination signal is received.
	// <-stop

	// // fmt.Println("Shutting down the node...")
	// // if err := node.Stop(); err != nil {
	// // 	log.Printf("Error during node shutdown: %v\n", err)
	// // }
	// fmt.Println("Node has been successfully stopped.")
}
