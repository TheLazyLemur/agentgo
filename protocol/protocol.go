package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type AcpConnection struct {
	provider  IOProvider
	reader    io.Reader
	writer    io.Writer
	sessionID string
	recorder  ConversationRecorder
}

// OpenAcpConnection creates a new ACP connection with the given IO provider
func OpenAcpConnection(provider IOProvider) (*AcpConnection, error) {
	if err := provider.Start(); err != nil {
		return nil, err
	}

	return &AcpConnection{
		provider: provider,
		reader:   provider.GetReader(),
		writer:   provider.GetWriter(),
	}, nil
}

// OpenAcpStdioConnection creates a new ACP connection using binary execution (backward compatible)
func OpenAcpStdioConnection(
	command string,
	args ...string,
) *AcpConnection {
	provider := NewBinaryIOProvider(command, args...)
	conn, err := OpenAcpConnection(provider)
	if err != nil {
		panic(err)
	}
	return conn
}

// OpenAcpRecordingConnection creates a new ACP connection with recording enabled
func OpenAcpRecordingConnection(command, recordingFile string, args ...string) (*AcpConnection, error) {
	provider := NewBinaryIOProvider(command, args...)
	conn, err := OpenAcpConnection(provider)
	if err != nil {
		return nil, err
	}

	// Add recorder
	recorder, err := NewFileRecorder(recordingFile)
	if err != nil {
		conn.Close() // Clean up the connection
		return nil, err
	}
	conn.recorder = recorder
	return conn, nil
}

// OpenAcpReplayConnection creates a connection that replays recorded messages
func OpenAcpReplayConnection(recordingFile string) (*AcpConnection, error) {
	provider, err := NewReplayIOProvider(recordingFile)
	if err != nil {
		return nil, err
	}
	return OpenAcpConnection(provider)
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

	_, err = acpConn.writer.Write(append(data, '\n'))
	if err != nil {
		return "", fmt.Errorf("failed to write request to gemini: %v", err)
	}

	response := map[string]any{}
	if err := json.NewDecoder(acpConn.reader).Decode(&response); err != nil {
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

// Close closes the connection and cleans up resources
func (acpConn *AcpConnection) Close() error {
	var err error

	// Close recorder first to save recording
	if acpConn.recorder != nil {
		if recErr := acpConn.recorder.Close(); recErr != nil {
			err = recErr
		}
	}

	// Close provider
	if acpConn.provider != nil {
		if provErr := acpConn.provider.Close(); provErr != nil && err == nil {
			err = provErr
		}
	}

	return err
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

	_, err = acpConn.writer.Write(append(data, '\n'))
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

	_, err = acpConn.writer.Write(append(data, '\n'))
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
	decoder := json.NewDecoder(acpConn.reader)
	
	for {
		response := map[string]any{}
		if err := decoder.Decode(&response); err != nil {
			return fmt.Errorf("JSON decode error in StreamResponses: %v", err)
		}

		// Record the response if recorder is present
		if acpConn.recorder != nil {
			if err := acpConn.recorder.RecordMessage(response); err != nil {
				return fmt.Errorf("failed to record message: %v", err)
			}
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
		// For recording, I think we should add a setting to the acpConnection to record request to file
		// This can then later be used for replays/
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
