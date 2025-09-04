package claude

import (
	"agentgo/protocol"
	"encoding/json"
	"fmt"
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
