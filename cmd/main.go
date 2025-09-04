package main

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

func main() {
	ch := make(chan int)

	acpConn := protocol.OpenAcpStdioConnection("claude-code-acp")
	_, err := acpConn.InitializeSession()
	if err != nil {
		panic(err)
	}

	claude := &claude.Claude{}
	handlers := handlers.NewHandlers(claude, claude)

	go func() {
		if err := acpConn.StreamResponses(handlers, ch); err != nil {
			panic(err)
		}
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
}
