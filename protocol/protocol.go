package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type AcpConnection struct {
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	sessionID string
}

func OpenAcpStdioConnection(
	command string,
	args ...string,
) *AcpConnection {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic("failed to create stdin pipe")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		panic("failed to create stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		panic("failed to start gemini process")
	}

	return &AcpConnection{
		stdin:  stdin,
		stdout: stdout,
	}
}

func (acpConn *AcpConnection) InitializeSession() (string, error) {
	cwd, _ := os.Getwd()
	sessionNewReq := SessionNewRequest{
		JSONRPC: "2.0",
		ID:      0,
		Method:  "session/new",
		Params: SessionParams{
			Cwd:        cwd,
			MCPServers: make([]MCPServer, 0),
		},
	}

	data, err := json.Marshal(sessionNewReq)
	if err != nil {
		return "", fmt.Errorf("failed to encode request to gemini: %v", err)
	}

	_, err = acpConn.stdin.Write(append(data, '\n'))
	if err != nil {
		return "", fmt.Errorf("failed to write request to gemini: %v", err)
	}

	response := map[string]any{}
	if err := json.NewDecoder(acpConn.stdout).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response from gemini: %v", err)
	}

	if _, ok := response["result"]; !ok {
		return "", fmt.Errorf("expected result, did not get it %v", response)
	}

	resultMap := response["result"].(map[string]any)

	if sessionID, hasSessionID := resultMap["sessionId"]; hasSessionID {
		if _, ok := sessionID.(string); !ok {
			return "", fmt.Errorf("expected sessionID to be type of string: %v", sessionID)
		}

		acpConn.sessionID = sessionID.(string)
		return sessionID.(string), nil
	}

	return "", fmt.Errorf("")
}

func (acpConn *AcpConnection) HasInit() bool {
	return acpConn.sessionID != ""
}

func (acpConn *AcpConnection) SendMessage(message string) error {
	promptReq := SessionPromptRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "session/prompt",
		Params: SessionPromptParams{
			SessionID: acpConn.sessionID,
			Prompt: []Prompt{
				{Type: "text", Text: message},
			},
		},
	}

	data, err := json.Marshal(promptReq)
	if err != nil {
		return err
	}

	_, err = acpConn.stdin.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	return nil
}

func (acpConn *AcpConnection) SendToolResponse(reqID int, optionID string) error {
	resp := ToolPermissionResponse{
		JSONRPC: "2.0",
		ID:      reqID,
		Result: ToolPermissionResult{
			Outcome: ToolPermissionOutcome{
				Outcome:  "selected",
				OptionID: optionID,
			},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = acpConn.stdin.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	return nil
}

type handler interface {
	HandlePermissionRequest(
		acpConn *AcpConnection,
		raw []byte,
		req SessionRequestPermissionRequest,
	) error
	HandleNotification(raw []byte, req SessionUpdateRequest) error
}

func (acpConn *AcpConnection) StreamResponses(handlers handler, ch chan int) error {
	for {
		response := map[string]any{}
		if err := json.NewDecoder(acpConn.stdout).Decode(&response); err != nil {
			return err
		}

		if err := Callback(handlers, ch)(acpConn, response); err != nil {
			return err
		}
	}
}

func Callback(
	handlers handler,
	ch chan int,
) func(acpConn *AcpConnection, response map[string]any) error {
	return func(acpConn *AcpConnection, response map[string]any) error {
		method, ok := response["method"].(string)
		if !ok || method == "" {
			ch <- 0
			return nil
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			return err
		}

		switch method {
		case "session/update":
			var req SessionUpdateRequest
			if err := json.Unmarshal(jsonData, &req); err != nil {
				return err
			}
			if err := handlers.HandleNotification(jsonData, req); err != nil {
				return err
			}
		case "session/request_permission":
			var req SessionRequestPermissionRequest
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
