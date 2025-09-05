package protocol

import (
	"io"
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