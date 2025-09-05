package main

import (
	"fmt"
	"os"

	"agentgo/internal/app"
)

func main() {
	// Create and run application coordinator
	coordinator, err := app.NewCoordinator()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer coordinator.Close()

	// Run the main application
	if err := coordinator.Run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
}
