package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"agentgo/protocol"
)

// LifecycleManager handles application startup and graceful shutdown
type LifecycleManager struct {
	connection *protocol.AcpConnection
	sigChan    chan os.Signal
}

// NewLifecycleManager creates a new lifecycle manager
func NewLifecycleManager(connection *protocol.AcpConnection) *LifecycleManager {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	return &LifecycleManager{
		connection: connection,
		sigChan:    sigChan,
	}
}

// SetupGracefulShutdown sets up graceful shutdown handling in a goroutine
func (lm *LifecycleManager) SetupGracefulShutdown() {
	go func() {
		<-lm.sigChan
		fmt.Println("\nShutting down gracefully...")
		if err := lm.connection.Close(); err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
		os.Exit(0)
	}()
}