package protocol

import (
	"testing"
)

func TestAcpConnection_Construction(t *testing.T) {
	// Test OpenAcpStdioConnection (this will fail due to missing binary, but interface should work)
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic due to missing claude-code-acp binary
				// This just tests that the interface works
			}
		}()
		_ = OpenAcpStdioConnection("non-existent-binary")
	}()

	_, err := OpenAcpRecordingConnection("non-existent-binary", "test.jsonl")
	if err == nil {
		t.Error("Expected error for non-existent binary")
	}

	_, err = OpenAcpReplayConnection("non-existent-file.jsonl")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestMessageRouting(t *testing.T) {
	ch := make(chan int, 1)

	handler := &MockHandler{}

	response := map[string]any{
		"method": "unknown/method",
	}

	err := RouteMessage(handler, ch, nil, response)
	if err != nil {
		t.Errorf("RouteMessage should not error on unknown methods: %v", err)
	}

	select {
	case <-ch:
		t.Error("Should not signal completion for empty method response")
	default:
	}
}

type MockHandler struct{}

func (m *MockHandler) HandlePermissionRequest(*AcpConnection, []byte, SessionRequestPermissionRequest) error {
	return nil
}

func (m *MockHandler) HandleNotification([]byte, SessionUpdateRequest) error {
	return nil
}
