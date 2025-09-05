package protocol

import (
	"os"
	"testing"
	"time"
)

func TestReplayConnection(t *testing.T) {
	// Create a sample recording file
	recordingFile := "test_recording.jsonl"
	defer os.Remove(recordingFile)

	testMessages := []ConversationMessage{
		{
			Timestamp: time.Now(),
			Data: map[string]any{
				"jsonrpc": "2.0",
				"result": map[string]any{
					"sessionId": "test123",
				},
			},
		},
		{
			Timestamp: time.Now(),
			Data: map[string]any{
				"method": "session/update",
				"params": map[string]any{
					"content": "test message",
				},
			},
		},
	}

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

	conn, err := OpenAcpReplayConnection(recordingFile)
	if err != nil {
		t.Fatalf("Failed to create replay connection: %v", err)
	}
	defer conn.Close()

	if conn.provider == nil {
		t.Fatal("Provider is nil")
	}

	if conn.reader == nil {
		t.Fatal("Reader is nil")
	}

	if conn.writer == nil {
		t.Fatal("Writer is nil")
	}
}
