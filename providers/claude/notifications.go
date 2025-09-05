package claude

import (
	"encoding/json"
	"fmt"
	"strings"

	"agentgo/protocol"
)

func (c *Claude) HandleNotification(raw []byte, req protocol.SessionUpdateRequest) error {
	var data map[string]any
	err := json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Printf("Error parsing notification: %v\n", err)
		return nil
	}

	params, ok := data["params"].(map[string]any)
	if !ok {
		return nil
	}
	update, ok := params["update"].(map[string]any)
	if !ok {
		return nil
	}

	updateType, _ := update["sessionUpdate"].(string)
	if updateType == "plan" {
		return c.handleTodoListUpdate(update)
	}

	content, ok := update["content"].(map[string]any)
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

func (c *Claude) handleTodoListUpdate(update map[string]any) error {
	entries, ok := update["entries"].([]any)
	if !ok {
		return nil
	}

	fmt.Println("\n\033[1;35m📋 Todo List Update:\033[0m")
	fmt.Println("\033[1;37m" + strings.Repeat("─", 50) + "\033[0m")

	for i, entry := range entries {
		item, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		content, _ := item["content"].(string)
		status, _ := item["status"].(string)
		priority, _ := item["priority"].(string)

		var statusIcon, statusColor string
		switch status {
		case "completed":
			statusIcon = "✅"
			statusColor = "\033[1;32m"
		case "in_progress":
			statusIcon = "🔄"
			statusColor = "\033[1;33m"
		case "pending":
			statusIcon = "⏳"
			statusColor = "\033[1;36m"
		default:
			statusIcon = "❓"
			statusColor = "\033[1;37m"
		}

		priorityIndicator := ""
		switch priority {
		case "high":
			priorityIndicator = " \033[1;31m[HIGH]\033[0m"
		case "low":
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
