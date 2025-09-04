package protocol

type InitializeRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	ProtocolVersion    int                `json:"protocolVersion"`
	ClientCapabilities ClientCapabilities `json:"clientCapabilities"`
}

type ClientCapabilities struct {
	FS FileSystemCapabilities `json:"fs"`
}

type FileSystemCapabilities struct {
	ReadTextFile  bool `json:"readTextFile"`
	WriteTextFile bool `json:"writeTextFile"`
}

type SessionNewRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  SessionParams `json:"params"`
}

type SessionParams struct {
	Cwd        string      `json:"cwd"`
	MCPServers []MCPServer `json:"mcpServers"`
}

type MCPServer struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Env     []string `json:"env"`
}

type SessionPromptRequest struct {
	JSONRPC string              `json:"jsonrpc"`
	ID      int                 `json:"id"`
	Method  string              `json:"method"`
	Params  SessionPromptParams `json:"params"`
}

type SessionPromptParams struct {
	SessionID string   `json:"sessionId"`
	Prompt    []Prompt `json:"prompt"`
}

type Prompt struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	Resource *PromptResource `json:"resource,omitempty"`
}

type PromptResource struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

type SessionUpdateRequest struct {
	JSONRPC string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  SessionUpdateParams `json:"params"`
}

type SessionUpdateParams struct {
	SessionID string        `json:"sessionId"`
	Update    SessionUpdate `json:"update"`
	Params    SessionUpdate `json:"params"`
}

type SessionUpdate struct {
	AvailableCommands []Command `json:"availableCommands,omitempty"`
	SessionUpdateType string    `json:"sessionUpdate"`
}

type Command struct {
	Description string `json:"description"`
	Input       any    `json:"input"`
	Name        string `json:"name"`
}

type SessionRequestPermissionRequest struct {
	JSONRPC string                         `json:"jsonrpc"`
	ID      int                            `json:"id"`
	Method  string                         `json:"method"`
	Params  SessionRequestPermissionParams `json:"params"`
}

type SessionRequestPermissionParams struct {
	SessionID string             `json:"sessionId"`
	ToolCall  ToolCall           `json:"toolCall"`
	Options   []PermissionOption `json:"options"`
}

const (
	OptionIDAllowAlways = "allow_always"
	OptionIDAllowOnce   = "allow"
	OptionIDRejectOnce  = "reject"
)

type ToolCall struct {
	ToolCallID string         `json:"toolCallId"`
	RawInput   map[string]any `json:"rawInput"`
}

type ToolPermissionResponse struct {
	JSONRPC string               `json:"jsonrpc"`
	ID      int                  `json:"id"`
	Result  ToolPermissionResult `json:"result"`
}

type ToolPermissionResult struct {
	Outcome ToolPermissionOutcome `json:"outcome"`
}

type ToolPermissionOutcome struct {
	Outcome  string `json:"outcome"`
	OptionID string `json:"optionId"`
}

type PermissionOption struct {
	OptionID string `json:"optionId"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
}

type Response struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  *ResponseResult `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

type ResponseResult struct {
	StopReason string `json:"stopReason,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
