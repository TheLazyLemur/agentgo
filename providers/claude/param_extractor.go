package claude

import (
	"encoding/json"
)

// ExtractToolParams returns the raw parameters
func ExtractToolParams(raw json.RawMessage) map[string]any {
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
