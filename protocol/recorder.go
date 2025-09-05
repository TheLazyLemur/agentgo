package protocol

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// ConversationMessage represents a single message in a recorded conversation
type ConversationMessage struct {
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// RecordedConversation represents a complete recorded conversation
type RecordedConversation struct {
	SessionID string                `json:"session_id,omitempty"`
	Messages  []ConversationMessage `json:"messages"`
}

// ConversationRecorder interface for recording conversations
type ConversationRecorder interface {
	RecordMessage(data map[string]interface{}) error
	Close() error
}

// FileRecorder implements ConversationRecorder by writing to a JSON file
type FileRecorder struct {
	filePath string
	file     *os.File
	encoder  *json.Encoder
	mutex    sync.Mutex
}

// NewFileRecorder creates a new file recorder
func NewFileRecorder(filePath string) (*FileRecorder, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	
	return &FileRecorder{
		filePath: filePath,
		file:     file,
		encoder:  json.NewEncoder(file),
	}, nil
}

// RecordMessage records a message to the conversation
func (f *FileRecorder) RecordMessage(data map[string]interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	message := ConversationMessage{
		Timestamp: time.Now(),
		Data:      data,
	}

	// Write each message as a separate JSON line
	return f.encoder.Encode(message)
}

// Close closes the file and flushes any remaining data
func (f *FileRecorder) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.file != nil {
		return f.file.Close()
	}
	return nil
}
