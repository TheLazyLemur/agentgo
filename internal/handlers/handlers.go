package handlers

import (
	"agentgo/protocol"
	"agentgo/internal/core"
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
	return h.toolPrompt.HandlePermissionRequest(acpConn, raw, req)
}

func (h *Handlers) HandleNotification(raw []byte, req protocol.SessionUpdateRequest) error {
	return h.notifier.HandleNotification(raw, req)
}
