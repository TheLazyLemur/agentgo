package claude

import (
	"testing"

	"agentgo/protocol"
)

func TestDisplayToolRequest(t *testing.T) {
	toolType := ToolBash
	toolID := "test-tool-123"
	params := map[string]any{
		"💻 Command": "ls -la",
		"old_string": "old content",
		"new_string": "new content",
	}
	options := []protocol.PermissionOption{
		{OptionID: "allow", Name: "Allow once"},
		{OptionID: "allow_always", Name: "Allow always"},
		{OptionID: "reject", Name: "Reject"},
	}

	// This test just verifies the function doesn't crash
	// In a real scenario, we might capture stdout to verify output
	err := DisplayToolRequest(toolType, toolID, params, options)
	if err != nil {
		t.Errorf("DisplayToolRequest() returned error: %v", err)
	}
}

func TestShowUserSelection(t *testing.T) {
	err := ShowUserSelection("Allow once")
	if err != nil {
		t.Errorf("ShowUserSelection() returned error: %v", err)
	}
}