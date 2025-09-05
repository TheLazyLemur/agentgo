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

func TestFormatParamsForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "bash command formatting",
			input: map[string]any{
				"command": "ls -la /home",
			},
			expected: map[string]any{
				"💻 Command": "ls -la /home",
			},
		},
		{
			name: "file operations formatting",
			input: map[string]any{
				"file_path": "/path/to/file.txt",
				"content":   "Hello world",
			},
			expected: map[string]any{
				"📁 File":    "/path/to/file.txt",
				"📝 Content": "Hello world",
			},
		},
		{
			name: "content truncation",
			input: map[string]any{
				"content": string(make([]rune, 250)) + "extra",
			},
			expected: map[string]any{
				"📝 Content": string(make([]rune, 200)) + "...",
			},
		},
		{
			name: "string truncation",
			input: map[string]any{
				"old_string": string(make([]rune, 150)),
			},
			expected: map[string]any{
				"old_string": string(make([]rune, 100)) + "...",
			},
		},
		{
			name: "edits count formatting",
			input: map[string]any{
				"edits": []any{
					map[string]any{"old": "1", "new": "1"},
					map[string]any{"old": "2", "new": "2"},
				},
			},
			expected: map[string]any{
				"📝 Edits": "2 edit(s)",
			},
		},
		{
			name: "passthrough other fields",
			input: map[string]any{
				"unknown_field": "value",
				"number":        42,
			},
			expected: map[string]any{
				"unknown_field": "value",
				"number":        42,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatParamsForDisplay(tt.input)
			
			for key, expectedValue := range tt.expected {
				if actualValue, ok := result[key]; !ok {
					t.Errorf("formatParamsForDisplay() missing key %q", key)
				} else if actualValue != expectedValue {
					t.Errorf("formatParamsForDisplay()[%q] = %v, expected %v", key, actualValue, expectedValue)
				}
			}
		})
	}
}