package claude

import (
	"fmt"
	"strings"

	"agentgo/protocol"
)

// DisplayToolRequest shows a formatted tool permission request to the user
func DisplayToolRequest(toolType ToolType, toolID string, rawParams map[string]any, options []protocol.PermissionOption) error {
	// Format raw parameters for Claude's distinctive display
	params := formatParamsForDisplay(rawParams)
	
	// Enhanced tool output formatting
	fmt.Printf("\n\033[1;36m╭─ Tool Request ─────────────────────────────────╮\033[0m\n")
	fmt.Printf("\033[1;36m│\033[0m \033[1;33m🔧 %s\033[0m\n", toolType)
	fmt.Printf("\033[1;36m│\033[0m ID: \033[0;37m%s\033[0m\n", toolID)

	// Display parameters in a more readable format
	if len(params) > 0 {
		fmt.Printf("\033[1;36m│\033[0m\n\033[1;36m│\033[0m \033[1;32mParameters:\033[0m\n")
		for key, value := range params {
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
	for i, option := range options {
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
	
	return nil
}

// PromptUserChoice asks the user to select from the available options
func PromptUserChoice(numOptions int) (int, error) {
	var choice int
	fmt.Printf("\n\033[1;33m❓ Select your choice (1-%d):\033[0m ", numOptions)
	fmt.Scanf("%d", &choice)

	if choice < 1 || choice > numOptions {
		choice = 2 // default to "allow" if invalid
	}
	
	return choice, nil
}

// ShowUserSelection displays the user's choice
func ShowUserSelection(selectedOption string) error {
	fmt.Printf("\033[1;32m✓ Selected:\033[0m %s\n\n", selectedOption)
	return nil
}

// formatParamsForDisplay formats raw parameters for Claude's distinctive UI display
func formatParamsForDisplay(rawParams map[string]any) map[string]any {
	enhanced := make(map[string]any)
	for key, value := range rawParams {
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