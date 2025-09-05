package app

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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

	connection, err := createConnection(config)
	if err != nil {
		return nil, err
	}

	if !config.IsReplaying() {
		_, err = connection.InitializeSession()
		if err != nil {
			return nil, err
		}
	}

	claude := &claude.Claude{}
	appHandlers := NewHandlers(claude, claude)

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

	c.lifecycle.SetupGracefulShutdown()

	go func() {
		if err := c.connection.StreamResponses(c.handlers, ch); err != nil {
			panic(err)
		}
	}()

	return c.runInteractionLoop(ch)
}

func (c *Coordinator) runInteractionLoop(ch chan int) error {
	for {
		time.Sleep(time.Second * 5)
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("> ")
		line, err := reader.ReadString('\n')
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
