package claude

import (
	"encoding/json"

	"agentgo/protocol"
)

// NotificationType represents the type of notification received
type NotificationType string

const (
	NotificationAgentChunk NotificationType = "agent_message_chunk"
	NotificationUser       NotificationType = "user_message"
	NotificationTodoList   NotificationType = "plan"
	NotificationGeneric    NotificationType = "generic"
)

// NotificationData holds parsed notification information
type NotificationData struct {
	Type        NotificationType
	Text        string
	ContentType string
	UpdateType  string
}

// TodoEntry represents a single todo list entry
type TodoEntry struct {
	Content  string
	Status   string
	Priority string
}

// TodoListData holds parsed todo list information
type TodoListData struct {
	Entries []TodoEntry
}

// ParseNotification extracts the core notification data from raw JSON
func ParseNotification(raw []byte, req protocol.SessionUpdateRequest) (*NotificationData, error) {
	var data map[string]any
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}

	params, ok := data["params"].(map[string]any)
	if !ok {
		return nil, nil
	}

	update, ok := params["update"].(map[string]any)
	if !ok {
		return nil, nil
	}

	updateType, _ := update["sessionUpdate"].(string)

	if updateType == "plan" {
		return &NotificationData{
			Type:       NotificationTodoList,
			UpdateType: updateType,
		}, nil
	}

	content, ok := update["content"].(map[string]any)
	if !ok {
		return nil, nil
	}

	text, _ := content["text"].(string)
	contentType, _ := content["type"].(string)

	if text == "" {
		return nil, nil
	}

	var notificationType NotificationType
	switch updateType {
	case "agent_message_chunk":
		notificationType = NotificationAgentChunk
	case "user_message":
		notificationType = NotificationUser
	default:
		notificationType = NotificationGeneric
	}

	return &NotificationData{
		Type:        notificationType,
		Text:        text,
		ContentType: contentType,
		UpdateType:  updateType,
	}, nil
}

// ParseTodoList extracts todo list data from the update map
func ParseTodoList(update map[string]any) *TodoListData {
	entries, ok := update["entries"].([]any)
	if !ok {
		return nil
	}

	var todoEntries []TodoEntry
	for _, entry := range entries {
		item, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		content, _ := item["content"].(string)
		status, _ := item["status"].(string)
		priority, _ := item["priority"].(string)

		todoEntries = append(todoEntries, TodoEntry{
			Content:  content,
			Status:   status,
			Priority: priority,
		})
	}

	return &TodoListData{
		Entries: todoEntries,
	}
}
