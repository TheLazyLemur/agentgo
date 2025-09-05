package claude

import (
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

	// Display the tool request to the user
	if err := DisplayToolRequest(toolType, req.Params.ToolCall.ToolCallID, toolParams, req.Params.Options); err != nil {
		return err
	}

	// Get user choice
	choice, err := PromptUserChoice(len(req.Params.Options))
	if err != nil {
		return err
	}

	selectedOption := req.Params.Options[choice-1]
	
	// Show the user's selection
	if err := ShowUserSelection(selectedOption.Name); err != nil {
		return err
	}

	return acpConn.SendToolResponse(req.ID, selectedOption.OptionID)
}
