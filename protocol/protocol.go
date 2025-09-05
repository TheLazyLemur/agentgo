package protocol

// This file exists for backward compatibility.
// All protocol functionality has been moved to more focused modules:
//
// - connection.go: AcpConnection struct and connection management
// - message_router.go: Message routing and handler interface
// - types.go: Protocol type definitions (already separate)
// - io_provider.go: I/O provider abstractions (already separate)
// - recorder.go: Conversation recording (already separate)
//
// Import protocol package to access all functionality.
