package protocol

import (
	"encoding/json"
	"fmt"
)

// Handler interface for handling different types of protocol messages
type Handler interface {
	HandlePermissionRequest(
		acpConn *AcpConnection,
		raw []byte,
		req SessionRequestPermissionRequest,
	) error
	HandleNotification(raw []byte, req SessionUpdateRequest) error
}

// StreamResponses processes incoming messages and routes them to appropriate handlers
func (acpConn *AcpConnection) StreamResponses(handlers Handler, ch chan int) error {
	decoder := json.NewDecoder(acpConn.reader)
	
	for {
		response := map[string]any{}
		if err := decoder.Decode(&response); err != nil {
			return fmt.Errorf("JSON decode error in StreamResponses: %v", err)
		}

		// Record the response if recorder is present
		if acpConn.recorder != nil {
			if err := acpConn.recorder.RecordMessage(response); err != nil {
				return fmt.Errorf("failed to record message: %v", err)
			}
		}

		if err := RouteMessage(handlers, ch, acpConn, response); err != nil {
			return err
		}
	}
}

// RouteMessage routes a single message to the appropriate handler based on method
func RouteMessage(
	handlers Handler,
	ch chan int,
	acpConn *AcpConnection,
	response map[string]any,
) error {
	// For recording, I think we should add a setting to the acpConnection to record request to file
	// This can then later be used for replays/
	method, ok := response["method"].(string)
	if !ok || method == "" {
		ch <- 0
		return nil
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}

	switch method {
	case "session/update":
		var req SessionUpdateRequest
		if err := json.Unmarshal(jsonData, &req); err != nil {
			return err
		}
		if err := handlers.HandleNotification(jsonData, req); err != nil {
			return err
		}
	case "session/request_permission":
		var req SessionRequestPermissionRequest
		if err := json.Unmarshal(jsonData, &req); err != nil {
			return err
		}
		if err := handlers.HandlePermissionRequest(acpConn, jsonData, req); err != nil {
			return err
		}
	default:
		fmt.Println()
		fmt.Println("======UNKNOWN=======")
		fmt.Println(string(jsonData))
		fmt.Println("=====================")
	}

	return nil
}

// Callback is a backward compatible function that creates a message routing callback
// Deprecated: Use RouteMessage directly
func Callback(
	handlers Handler,
	ch chan int,
) func(acpConn *AcpConnection, response map[string]any) error {
	return func(acpConn *AcpConnection, response map[string]any) error {
		return RouteMessage(handlers, ch, acpConn, response)
	}
}