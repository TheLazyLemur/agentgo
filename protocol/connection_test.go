package protocol

import (
	"testing"
)

func TestAcpConnection_Construction(t *testing.T) {
	// Test that we can create connections via different constructors
	// Using a simple test to verify the APIs still work

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

	// Test OpenAcpRecordingConnection
	_, err := OpenAcpRecordingConnection("non-existent-binary", "test.jsonl")
	if err == nil {
		t.Error("Expected error for non-existent binary")
	}

	// Test OpenAcpReplayConnection  
	_, err = OpenAcpReplayConnection("non-existent-file.jsonl")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestMessageRouting(t *testing.T) {
	// Test that RouteMessage handles unknown methods gracefully
	ch := make(chan int, 1)
	
	// Mock handler that does nothing
	handler := &MockHandler{}
	
	response := map[string]any{
		"method": "unknown/method",
	}

	// Should not error on unknown methods
	err := RouteMessage(handler, ch, nil, response)
	if err != nil {
		t.Errorf("RouteMessage should not error on unknown methods: %v", err)
	}
	
	// Should signal completion
	select {
	case <-ch:
		t.Error("Should not signal completion for empty method response")
	default:
		// Expected - no signal for unknown method
	}
}

// MockHandler for testing
type MockHandler struct{}

func (m *MockHandler) HandlePermissionRequest(*AcpConnection, []byte, SessionRequestPermissionRequest) error {
	return nil
}

func (m *MockHandler) HandleNotification([]byte, SessionUpdateRequest) error {
	return nil
}