package claude

import (
	"encoding/json"
	"testing"
)

func TestDetectToolType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ToolType
	}{
		{
			name: "bash tool",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"command": "ls -la"
						}
					}
				}
			}`,
			expected: ToolBash,
		},
		{
			name: "write tool",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"file_path": "/path/to/file.txt",
							"content": "Hello world"
						}
					}
				}
			}`,
			expected: ToolWrite,
		},
		{
			name: "edit tool",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"file_path": "/path/to/file.txt",
							"old_string": "old content",
							"new_string": "new content"
						}
					}
				}
			}`,
			expected: ToolEdit,
		},
		{
			name: "multi-edit tool",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"file_path": "/path/to/file.txt",
							"edits": [
								{"old_string": "old1", "new_string": "new1"},
								{"old_string": "old2", "new_string": "new2"}
							]
						}
					}
				}
			}`,
			expected: ToolMultiEdit,
		},
		{
			name: "unknown tool",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"unknown_field": "value"
						}
					}
				}
			}`,
			expected: ToolUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := json.RawMessage(tt.input)
			result := DetectToolType(raw)
			if result != tt.expected {
				t.Errorf("DetectToolType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}