package handlers

import (
	"agentgo/internal/core"
	"agentgo/protocol"
)

type Handlers struct {
	toolPrompt core.PermissionPrompt
	notifier   core.Notifications
}

func NewHandlers(toolPrompt core.PermissionPrompt, notifier core.Notifications) *Handlers {
	return &Handlers{
		toolPrompt: toolPrompt,
		notifier:   notifier,
	}
}

func (h *Handlers) HandlePermissionRequest(
	acpConn *protocol.AcpConnection,
	raw []byte,
	req protocol.SessionRequestPermissionRequest,
) error {
	// fmt.Println()
	// fmt.Println("======SESSION/REQUEST_PERMISSION=======")
	// fmt.Println(string(raw))
	// fmt.Println("=====================")

	return h.toolPrompt.HandlePermissionRequest(acpConn, raw, req)
}

func (h *Handlers) HandleNotification(raw []byte, req protocol.SessionUpdateRequest) error {
	// fmt.Println()
	// fmt.Println("======SESSION/UPDATE=======")
	// fmt.Println(string(raw))
	// fmt.Println("=====================")

	return h.notifier.HandleNotification(raw, req)
}
