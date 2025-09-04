package claude

import (
	"agentgo/protocol"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Claude) HandleNotification(raw []byte, req protocol.SessionUpdateRequest) error {
	var data map[string]interface{}
	err := json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Printf("Error parsing notification: %v\n", err)
		return nil
	}

	params, ok := data["params"].(map[string]interface{})
	if !ok {
		return nil
	}

	update, ok := params["update"].(map[string]interface{})
	if !ok {
		return nil
	}

	updateType, _ := update["sessionUpdate"].(string)

	// Check if this is a todo list update (plan)
	if updateType == "plan" {
		return c.handleTodoListUpdate(update)
	}

	content, ok := update["content"].(map[string]interface{})
	if !ok {
		return nil
	}

	text, _ := content["text"].(string)
	contentType, _ := content["type"].(string)

	if text == "" {
		return nil
	}

	switch updateType {
	case "agent_message_chunk":
		fmt.Printf("\033[1;34m🤖 Assistant:\033[0m %s", text)
	case "user_message":
		fmt.Printf("\033[1;32m👤 You:\033[0m %s", text)
	default:
		if contentType == "text" {
			fmt.Printf("\033[1;37m💬 Message:\033[0m %s", text)
		} else {
			fmt.Printf("\033[1;33m📄 %s:\033[0m %s", updateType, text)
		}
	}

	fmt.Println()
	return nil
}

func (c *Claude) handleTodoListUpdate(update map[string]interface{}) error {
	entries, ok := update["entries"].([]interface{})
	if !ok {
		return nil
	}

	fmt.Println("\n\033[1;35m📋 Todo List Update:\033[0m")
	fmt.Println("\033[1;37m" + strings.Repeat("─", 50) + "\033[0m")

	for i, entry := range entries {
		item, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		content, _ := item["content"].(string)
		status, _ := item["status"].(string)
		priority, _ := item["priority"].(string)

		// Status emoji and color
		var statusIcon, statusColor string
		switch status {
		case "completed":
			statusIcon = "✅"
			statusColor = "\033[1;32m" // green
		case "in_progress":
			statusIcon = "🔄"
			statusColor = "\033[1;33m" // yellow
		case "pending":
			statusIcon = "⏳"
			statusColor = "\033[1;36m" // cyan
		default:
			statusIcon = "❓"
			statusColor = "\033[1;37m" // white
		}

		// Priority indicator
		priorityIndicator := ""
		if priority == "high" {
			priorityIndicator = " \033[1;31m[HIGH]\033[0m"
		} else if priority == "low" {
			priorityIndicator = " \033[1;34m[LOW]\033[0m"
		}

		fmt.Printf("%s%d. %s %s%s%s\033[0m\n",
			statusColor,
			i+1,
			statusIcon,
			content,
			priorityIndicator,
			"")
	}

	fmt.Println("\033[1;37m" + strings.Repeat("─", 50) + "\033[0m\n")
	return nil
}
