package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"agentgo/internal/handlers"
	"agentgo/protocol"
	"agentgo/providers/claude"
)

func callback(
	handlers *handlers.Handlers,
	ch chan int,
) func(acpConn *protocol.AcpConnection, response map[string]any) error {
	return func(acpConn *protocol.AcpConnection, response map[string]any) error {
		method, ok := response["method"].(string)
		if !ok || method == "" {
			// fmt.Println(response)
			ch <- 0
			return nil
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			return err
		}

		switch method {
		case "session/update":
			// fmt.Println()
			// fmt.Println("======SESSION/UPDATE=======")
			// fmt.Println(string(jsonData))
			// fmt.Println("=====================")
			var req protocol.SessionUpdateRequest
			if err := json.Unmarshal(jsonData, &req); err != nil {
				return err
			}
			if err := handlers.HandleNotification(jsonData, req); err != nil {
				return err
			}
		case "session/request_permission":
			// fmt.Println()
			// fmt.Println("======SESSION/REQUEST_PERMISSION=======")
			// fmt.Println(string(jsonData))
			// fmt.Println("=====================")
			var req protocol.SessionRequestPermissionRequest
			if err := json.Unmarshal(jsonData, &req); err != nil {
				return err
			}
			if err := handlers.HandlePermissionRequest(acpConn, jsonData, req); err != nil {
				return err
			}
		default:
			fmt.Println()
			fmt.Println("======UNKNOWN=======")
			fmt.Println(string(jsonData))
			fmt.Println("=====================")
		}

		return nil
	}
}

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
		if err := acpConn.StreamResponses(callback(handlers, ch)); err != nil {
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
