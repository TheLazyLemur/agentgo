package claude

import (
	"encoding/json"
)

// ExtractToolParams returns the raw parameters without any formatting.
// Pure data extraction - no UI concerns.
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

	// Return pure raw parameters without any formatting
	return ri
}

// ExtractRawParams is an alias for ExtractToolParams for backward compatibility.
// Both functions now return pure raw data without formatting.
func ExtractRawParams(raw json.RawMessage) map[string]any {
	return ExtractToolParams(raw)
}
