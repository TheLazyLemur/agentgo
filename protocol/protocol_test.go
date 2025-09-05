package protocol

import (
	"os"
	"testing"
	"time"
)

func TestReplayConnection(t *testing.T) {
	// Create a sample recording file
	recordingFile := "test_recording.jsonl"
	defer os.Remove(recordingFile) // Clean up after test

	// Create test messages
	testMessages := []ConversationMessage{
		{
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"jsonrpc": "2.0",
				"result": map[string]interface{}{
					"sessionId": "test123",
				},
			},
		},
		{
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"method": "session/update",
				"params": map[string]interface{}{
					"content": "test message",
				},
			},
		},
	}

	// Write test recording
	recorder, err := NewFileRecorder(recordingFile)
	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}

	for _, msg := range testMessages {
		if err := recorder.RecordMessage(msg.Data); err != nil {
			t.Fatalf("Failed to record message: %v", err)
		}
	}

	if err := recorder.Close(); err != nil {
		t.Fatalf("Failed to close recorder: %v", err)
	}

	// Test replay connection
	conn, err := OpenAcpReplayConnection(recordingFile)
	if err != nil {
		t.Fatalf("Failed to create replay connection: %v", err)
	}
	defer conn.Close()

	// Test that we got a valid connection with a provider
	if conn.provider == nil {
		t.Fatal("Provider is nil")
	}

	// Test that the reader and writer are accessible through the provider
	if conn.reader == nil {
		t.Fatal("Reader is nil")
	}

	if conn.writer == nil {
		t.Fatal("Writer is nil")
	}
}