package claude

import (
	"encoding/json"
	"fmt"
)

// ExtractToolParams returns the rawInput params with enhanced formatting for display.
// Accepts either the full RPC envelope or the rawInput object itself.
func ExtractToolParams(raw json.RawMessage) map[string]any {
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

// ExtractRawParams returns the raw parameters without enhanced formatting.
// Useful for business logic that doesn't need display formatting.
func ExtractRawParams(raw json.RawMessage) map[string]any {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return map[string]any{"error": err.Error()}
	}
	
	ri := extractRawInput(m)
	if ri == nil {
		return map[string]any{}
	}
	
	return ri
}