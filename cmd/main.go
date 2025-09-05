package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agentgo/internal/handlers"
	"agentgo/protocol"
	"agentgo/providers/claude"
)

func main() {
	// Parse command line flags
	recordFile := flag.String("record", "", "Record conversation to file")
	flag.Parse()

	ch := make(chan int)

	// Create ACP connection with or without recording
	var acpConn *protocol.AcpConnection
	var err error

	if *recordFile != "" {
		fmt.Printf("Recording conversation to: %s\n", *recordFile)
		acpConn, err = protocol.OpenAcpRecordingConnection("claude-code-acp", *recordFile)
		if err != nil {
			panic(err)
		}
	} else {
		acpConn = protocol.OpenAcpStdioConnection("claude-code-acp")
	}
	_, err = acpConn.InitializeSession()
	if err != nil {
		panic(err)
	}

	claude := &claude.Claude{}
	handlers := handlers.NewHandlers(claude, claude)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := acpConn.StreamResponses(handlers, ch); err != nil {
			panic(err)
		}
	}()

	// Handle graceful shutdown
	go func() {
		<-sigChan
		fmt.Println("\nShutting down gracefully...")
		if err := acpConn.Close(); err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
		os.Exit(0)
	}()

	for {
		if acpConn.HasInit() == false {
			continue
		}
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
		if err := acpConn.SendMessage(line); err != nil {
			panic(err)
		}
		_ = <-ch
	}

	// Clean up resources
	if err := acpConn.Close(); err != nil {
		fmt.Printf("Error closing connection: %v\n", err)
	}
}
