package claude

import (
	"agentgo/protocol"
	"reflect"
	"testing"
)

func TestParseNotification(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *NotificationData
	}{
		{
			name: "agent message chunk",
			input: `{
				"params": {
					"update": {
						"sessionUpdate": "agent_message_chunk",
						"content": {
							"text": "Hello world",
							"type": "text"
						}
					}
				}
			}`,
			expected: &NotificationData{
				Type:        NotificationAgentChunk,
				Text:        "Hello world",
				ContentType: "text",
				UpdateType:  "agent_message_chunk",
			},
		},
		{
			name: "user message",
			input: `{
				"params": {
					"update": {
						"sessionUpdate": "user_message",
						"content": {
							"text": "Hi there",
							"type": "text"
						}
					}
				}
			}`,
			expected: &NotificationData{
				Type:        NotificationUser,
				Text:        "Hi there",
				ContentType: "text",
				UpdateType:  "user_message",
			},
		},
		{
			name: "todo list notification",
			input: `{
				"params": {
					"update": {
						"sessionUpdate": "plan",
						"entries": []
					}
				}
			}`,
			expected: &NotificationData{
				Type:       NotificationTodoList,
				UpdateType: "plan",
			},
		},
		{
			name: "generic notification",
			input: `{
				"params": {
					"update": {
						"sessionUpdate": "custom_update",
						"content": {
							"text": "Custom message",
							"type": "custom"
						}
					}
				}
			}`,
			expected: &NotificationData{
				Type:        NotificationGeneric,
				Text:        "Custom message",
				ContentType: "custom",
				UpdateType:  "custom_update",
			},
		},
		{
			name:     "empty text should return nil",
			input: `{
				"params": {
					"update": {
						"sessionUpdate": "agent_message_chunk",
						"content": {
							"text": "",
							"type": "text"
						}
					}
				}
			}`,
			expected: nil,
		},
		{
			name:     "malformed json should return error",
			input:    `{"invalid": json}`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			raw := []byte(tt.input)
			req := protocol.SessionUpdateRequest{}

			// when
			result, err := ParseNotification(raw, req)

			// then
			if tt.name == "malformed json should return error" {
				if err == nil {
					t.Error("ParseNotification() expected error for malformed JSON, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseNotification() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseNotification() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestParseTodoList(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected *TodoListData
	}{
		{
			name: "valid todo list",
			input: map[string]any{
				"entries": []any{
					map[string]any{
						"content":  "Task 1",
						"status":   "pending",
						"priority": "high",
					},
					map[string]any{
						"content":  "Task 2", 
						"status":   "completed",
						"priority": "low",
					},
				},
			},
			expected: &TodoListData{
				Entries: []TodoEntry{
					{
						Content:  "Task 1",
						Status:   "pending",
						Priority: "high",
					},
					{
						Content:  "Task 2",
						Status:   "completed", 
						Priority: "low",
					},
				},
			},
		},
		{
			name: "empty entries",
			input: map[string]any{
				"entries": []any{},
			},
			expected: &TodoListData{
				Entries: nil,
			},
		},
		{
			name:     "no entries field",
			input:    map[string]any{},
			expected: nil,
		},
		{
			name: "invalid entry format",
			input: map[string]any{
				"entries": []any{
					"invalid_entry",
					map[string]any{
						"content": "Valid task",
						"status":  "pending",
					},
				},
			},
			expected: &TodoListData{
				Entries: []TodoEntry{
					{
						Content: "Valid task",
						Status:  "pending",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			// ... input provided

			// when
			result := ParseTodoList(tt.input)

			// then
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseTodoList() = %v, expected %v", result, tt.expected)
			}
		})
	}
}