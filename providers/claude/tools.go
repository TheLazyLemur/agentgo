package claude

import (
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
	toolParams := ExtractToolParams(raw)

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
