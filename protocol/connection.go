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

// InitializeSession initializes a new session and returns the session ID
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

// SendMessage sends a user message to the session
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

// SendToolResponse sends a tool permission response
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