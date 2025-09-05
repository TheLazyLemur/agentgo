package claude

import (
	"encoding/json"
	"fmt"
	"strings"

	"agentgo/protocol"
)

type Claude struct{}

func (c *Claude) HandlePermissionRequest(
	acpConn *protocol.AcpConnection,
	raw []byte,
	req protocol.SessionRequestPermissionRequest,
) error {
	toolType := DetectToolType(raw)
	toolParams := c.ExtractToolParams(raw)

	// Enhanced tool output formatting
	fmt.Printf("\n\033[1;36m╭─ Tool Request ─────────────────────────────────╮\033[0m\n")
	fmt.Printf("\033[1;36m│\033[0m \033[1;33m🔧 %s\033[0m\n", toolType)
	fmt.Printf("\033[1;36m│\033[0m ID: \033[0;37m%s\033[0m\n", req.Params.ToolCall.ToolCallID)

	// Display parameters in a more readable format
	if len(toolParams) > 0 {
		fmt.Printf("\033[1;36m│\033[0m\n\033[1;36m│\033[0m \033[1;32mParameters:\033[0m\n")
		for key, value := range toolParams {
			if strings.HasPrefix(key, "📁") || strings.HasPrefix(key, "💻") ||
				strings.HasPrefix(key, "📝") {
				fmt.Printf("\033[1;36m│\033[0m   %s: \033[0;33m%v\033[0m\n", key, value)
			} else {
				switch key {
				case "old_string":
					fmt.Printf("\033[1;36m│\033[0m   🔍 Replace: \033[0;31m%v\033[0m\n", value)
				case "new_string":
					fmt.Printf("\033[1;36m│\033[0m   ✏️  With: \033[0;32m%v\033[0m\n", value)
				default:
					fmt.Printf("\033[1;36m│\033[0m   • %s: \033[0;37m%v\033[0m\n", key, value)
				}
			}
		}
	}

	fmt.Printf("\033[1;36m│\033[0m\n\033[1;36m│\033[0m \033[1;32mOptions:\033[0m\n")
	for i, option := range req.Params.Options {
		var icon string
		switch option.OptionID {
		case "allow_always":
			icon = "✅"
		case "allow":
			icon = "👍"
		case "reject":
			icon = "❌"
		default:
			icon = "⚪"
		}
		fmt.Printf("\033[1;36m│\033[0m   [%d] %s %s\n", i+1, icon, option.Name)
	}
	fmt.Printf("\033[1;36m╰────────────────────────────────────────────────╯\033[0m\n")

	var choice int
	fmt.Printf("\n\033[1;33m❓ Select your choice (1-%d):\033[0m ", len(req.Params.Options))
	fmt.Scanf("%d", &choice)

	if choice < 1 || choice > len(req.Params.Options) {
		choice = 2 // default to "allow" if invalid
	}

	selectedOption := req.Params.Options[choice-1]
	fmt.Printf("\033[1;32m✓ Selected:\033[0m %s\n\n", selectedOption.Name)

	return acpConn.SendToolResponse(req.ID, selectedOption.OptionID)
}


// ExtractToolParamsClaude returns the rawInput params with enhanced formatting for display.
// Accepts either the full RPC envelope or the rawInput object itself.
func (c *Claude) ExtractToolParams(raw json.RawMessage) map[string]any {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return map[string]any{"error": err.Error()}
	}
	ri := extractRawInput(m)
	if ri == nil {
		return map[string]any{}
	}

	// Enhanced parameter formatting for better display
	enhanced := make(map[string]any)
	for key, value := range ri {
		switch key {
		case "file_path":
			if str, ok := value.(string); ok {
				enhanced["📁 File"] = str
			}
		case "command":
			if str, ok := value.(string); ok {
				enhanced["💻 Command"] = str
			}
		case "content":
			if str, ok := value.(string); ok {
				if len(str) > 200 {
					enhanced["📝 Content"] = str[:200] + "..."
				} else {
					enhanced["📝 Content"] = str
				}
			}
		case "old_string", "new_string":
			if str, ok := value.(string); ok {
				if len(str) > 100 {
					enhanced[key] = str[:100] + "..."
				} else {
					enhanced[key] = str
				}
			}
		case "edits":
			if edits, ok := value.([]any); ok {
				enhanced["📝 Edits"] = fmt.Sprintf("%d edit(s)", len(edits))
			}
		default:
			enhanced[key] = value
		}
	}
	return enhanced
}

