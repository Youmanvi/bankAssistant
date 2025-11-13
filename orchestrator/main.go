package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	config := LoadConfig()

	// Validate required configuration
	if config.RetellAPIKey == "" {
		log.Fatalf("RETELL_API_KEY environment variable is required")
	}

	// Create and start orchestrator
	orchestrator := NewOrchestrator(config)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := orchestrator.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan

	// Graceful shutdown
	if err := orchestrator.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Printf("Orchestrator stopped")
}
