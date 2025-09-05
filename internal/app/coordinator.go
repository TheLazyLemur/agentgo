package app

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"agentgo/internal/handlers"
	"agentgo/protocol"
	"agentgo/providers/claude"
)

// Coordinator orchestrates the main application flow
type Coordinator struct {
	config     *Config
	connection *protocol.AcpConnection
	lifecycle  *LifecycleManager
	handlers   protocol.Handler
}

// NewCoordinator creates a new application coordinator
func NewCoordinator() (*Coordinator, error) {
	config := ParseFlags()
	
	// Create connection based on configuration
	connection, err := createConnection(config)
	if err != nil {
		return nil, err
	}
	
	// Initialize session for non-replay connections
	if !config.IsReplaying() {
		_, err = connection.InitializeSession()
		if err != nil {
			return nil, err
		}
	}
	
	// Create handlers
	claude := &claude.Claude{}
	appHandlers := handlers.NewHandlers(claude, claude)
	
	// Create lifecycle manager
	lifecycle := NewLifecycleManager(connection)
	
	return &Coordinator{
		config:     config,
		connection: connection,
		lifecycle:  lifecycle,
		handlers:   appHandlers,
	}, nil
}

// Run starts the main application loop
func (c *Coordinator) Run() error {
	ch := make(chan int)
	
	// Setup graceful shutdown
	c.lifecycle.SetupGracefulShutdown()
	
	// Start message processing
	go func() {
		if err := c.connection.StreamResponses(c.handlers, ch); err != nil {
			panic(err)
		}
	}()
	
	// Main interaction loop
	return c.runInteractionLoop(ch)
}

// runInteractionLoop handles the main user interaction loop
func (c *Coordinator) runInteractionLoop(ch chan int) error {
	for {
		time.Sleep(time.Second * 5)
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("> ")
		line, err := reader.ReadString('\n') // reads until Enter
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		
		if err := c.connection.SendMessage(line); err != nil {
			return err
		}
		_ = <-ch
	}
	
	return nil
}

// Close cleans up resources
func (c *Coordinator) Close() error {
	if c.connection != nil {
		return c.connection.Close()
	}
	return nil
}

// createConnection creates the appropriate connection based on config
func createConnection(config *Config) (*protocol.AcpConnection, error) {
	switch {
	case config.IsReplaying():
		fmt.Printf("Replaying conversation from: %s\n", config.ReplayFile)
		return protocol.OpenAcpReplayConnection(config.ReplayFile)
	case config.IsRecording():
		fmt.Printf("Recording conversation to: %s\n", config.RecordFile)
		return protocol.OpenAcpRecordingConnection("claude-code-acp", config.RecordFile)
	default:
		return protocol.OpenAcpStdioConnection("claude-code-acp"), nil
	}
}