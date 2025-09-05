package claude

import "encoding/json"

type ToolType string

const (
	ToolWrite     ToolType = "write"
	ToolEdit      ToolType = "edit"
	ToolMultiEdit ToolType = "edit"
	ToolBash      ToolType = "bash"
	ToolUnknown   ToolType = "unknown"
)

// DetectToolType classifies Claude's rawInput payloads.
func DetectToolType(raw json.RawMessage) ToolType {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ToolUnknown
	}

	ri := extractRawInput(m)
	if ri == nil {
		return ToolUnknown
	}

	_, hasCmd := ri["command"]
	_, hasFile := ri["file_path"]
	_, hasContent := ri["content"]
	_, hasOld := ri["old_string"]
	_, hasNew := ri["new_string"]

	if editsVal, ok := ri["edits"]; ok {
		if edits, ok := editsVal.([]any); ok && hasFile && len(edits) > 0 {
			return ToolMultiEdit
		}
	}

	switch {
	case hasCmd:
		return ToolBash
	case hasFile && hasOld && hasNew:
		return ToolEdit
	case hasFile && hasContent:
		return ToolWrite
	default:
		return ToolUnknown
	}
}

func extractRawInput(m map[string]any) map[string]any {
	if params, ok := m["params"].(map[string]any); ok {
		if tc, ok := params["toolCall"].(map[string]any); ok {
			if ri, ok := tc["rawInput"].(map[string]any); ok {
				return ri
			}
		}
	}
	if m["file_path"] != nil || m["command"] != nil || m["content"] != nil || m["edits"] != nil ||
		m["old_string"] != nil ||
		m["new_string"] != nil {
		return m
	}
	return nil
}
