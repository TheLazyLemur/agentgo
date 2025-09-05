package claude

import (
	"encoding/json"
	"fmt"

	"agentgo/protocol"
)

// HandleNotification processes notification messages with Claude's distinctive UI
func (c *Claude) HandleNotification(raw []byte, req protocol.SessionUpdateRequest) error {
	notificationData, err := ParseNotification(raw, req)
	if err != nil {
		fmt.Printf("Error parsing notification: %v\n", err)
		return nil
	}

	if notificationData == nil {
		return nil
	}

	if notificationData.Type == NotificationTodoList {
		return c.handleTodoListUpdate(raw)
	}

	return DisplayNotification(notificationData)
}

func (c *Claude) handleTodoListUpdate(raw []byte) error {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	params, ok := data["params"].(map[string]any)
	if !ok {
		return nil
	}

	update, ok := params["update"].(map[string]any)
	if !ok {
		return nil
	}

	todoData := ParseTodoList(update)

	return DisplayTodoList(todoData)
}
