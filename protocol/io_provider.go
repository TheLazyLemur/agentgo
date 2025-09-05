package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

// IOProvider abstracts the I/O operations for ACP communication
type IOProvider interface {
	GetReader() io.Reader
	GetWriter() io.Writer
	Start() error
	Close() error
}

// BinaryIOProvider implements IOProvider for real binary execution
type BinaryIOProvider struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// NewBinaryIOProvider creates a new binary IO provider
func NewBinaryIOProvider(command string, args ...string) *BinaryIOProvider {
	return &BinaryIOProvider{
		cmd: exec.Command(command, args...),
	}
}

// Start initializes the binary and pipes
func (b *BinaryIOProvider) Start() error {
	var err error
	
	b.stdin, err = b.cmd.StdinPipe()
	if err != nil {
		return err
	}
	
	b.stdout, err = b.cmd.StdoutPipe()
	if err != nil {
		b.stdin.Close()
		return err
	}
	
	if err := b.cmd.Start(); err != nil {
		b.stdin.Close()
		b.stdout.Close()
		return err
	}
	
	return nil
}

// GetReader returns the stdout reader
func (b *BinaryIOProvider) GetReader() io.Reader {
	return b.stdout
}

// GetWriter returns the stdin writer
func (b *BinaryIOProvider) GetWriter() io.Writer {
	return b.stdin
}

// Close closes all pipes and waits for process to exit
func (b *BinaryIOProvider) Close() error {
	if b.stdin != nil {
		b.stdin.Close()
	}
	if b.stdout != nil {
		b.stdout.Close()
	}
	if b.cmd != nil && b.cmd.Process != nil {
		return b.cmd.Wait()
	}
	return nil
}

// ReplayIOProvider streams recorded messages for testing
type ReplayIOProvider struct {
	reader *bytes.Reader
}

// NewReplayIOProvider creates a replay provider from a JSONL recording file
func NewReplayIOProvider(recordingFile string) (*ReplayIOProvider, error) {
	file, err := os.Open(recordingFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		
		var raw map[string]json.RawMessage
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}
		
		if dataField, exists := raw["data"]; exists {
			buffer.Write(dataField)
			buffer.WriteByte('\n')
		}
	}

	return &ReplayIOProvider{
		reader: bytes.NewReader(buffer.Bytes()),
	}, scanner.Err()
}

// Start is a no-op for replay
func (r *ReplayIOProvider) Start() error {
	return nil
}

// GetReader returns the message stream
func (r *ReplayIOProvider) GetReader() io.Reader {
	return r.reader
}

// GetWriter returns a no-op writer
func (r *ReplayIOProvider) GetWriter() io.Writer {
	return &noopWriter{}
}

// Close is a no-op for replay
func (r *ReplayIOProvider) Close() error {
	return nil
}

// noopWriter discards all writes
type noopWriter struct{}

func (n *noopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}