package core

import "agentgo/protocol"

type PermissionPrompt interface {
	HandlePermissionRequest(
		acpConn *protocol.AcpConnection,
		raw []byte,
		req protocol.SessionRequestPermissionRequest,
	) error
}

type Notifications interface {
	HandleNotification(raw []byte, req protocol.SessionUpdateRequest) error
}
