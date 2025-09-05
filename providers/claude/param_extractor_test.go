package claude

import (
	"encoding/json"
	"testing"
)

func TestExtractToolParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]any
	}{
		{
			name: "bash command",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"command": "ls -la /home"
						}
					}
				}
			}`,
			expected: map[string]any{
				"command": "ls -la /home",
			},
		},
		{
			name: "file write with content",
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
			expected: map[string]any{
				"file_path": "/path/to/file.txt",
				"content":   "Hello world",
			},
		},
		{
			name: "edit with normal strings",
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
			expected: map[string]any{
				"file_path":  "/path/to/file.txt",
				"old_string": "old content",
				"new_string": "new content",
			},
		},
		{
			name: "multi-edit",
			input: `{
				"params": {
					"toolCall": {
						"rawInput": {
							"file_path": "/path/to/file.txt",
							"edits": [
								{"old_string": "old1", "new_string": "new1"},
								{"old_string": "old2", "new_string": "new2"},
								{"old_string": "old3", "new_string": "new3"}
							]
						}
					}
				}
			}`,
			expected: map[string]any{
				"file_path": "/path/to/file.txt",
				"edits": []any{
					map[string]any{"old_string": "old1", "new_string": "new1"},
					map[string]any{"old_string": "old2", "new_string": "new2"},
					map[string]any{"old_string": "old3", "new_string": "new3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := json.RawMessage(tt.input)
			result := ExtractToolParams(raw)

			for key := range tt.expected {
				if _, ok := result[key]; !ok {
					t.Errorf("ExtractToolParams() missing key %q", key)
				}
			}

			if tt.name == "multi-edit" {
				if edits, ok := result["edits"]; ok {
					if editsSlice, ok := edits.([]any); ok {
						if len(editsSlice) != 3 {
							t.Errorf("Expected 3 edits, got %d", len(editsSlice))
						}
					} else {
						t.Errorf("Expected edits to be []interface{}, got %T", edits)
					}
				}
			} else {
				for key, expectedValue := range tt.expected {
					if actualValue, ok := result[key]; !ok {
						t.Errorf("ExtractToolParams() missing key %q", key)
					} else if actualValue != expectedValue {
						t.Errorf("ExtractToolParams()[%q] = %v, expected %v", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

func TestExtractRawParams(t *testing.T) {
	input := `{
		"params": {
			"toolCall": {
				"rawInput": {
					"file_path": "/path/to/file.txt",
					"command": "ls -la",
					"content": "Hello world"
				}
			}
		}
	}`

	raw := json.RawMessage(input)
	result := ExtractToolParams(raw)

	expected := map[string]any{
		"file_path": "/path/to/file.txt",
		"command":   "ls -la",
		"content":   "Hello world",
	}

	for key, expectedValue := range expected {
		if actualValue, ok := result[key]; !ok {
			t.Errorf("ExtractRawParams() missing key %q", key)
		} else if actualValue != expectedValue {
			t.Errorf("ExtractRawParams()[%q] = %v, expected %v", key, actualValue, expectedValue)
		}
	}
}
