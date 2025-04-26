package log

import (
	"fmt"
	"os"
	"sync"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: f}, nil
}

func (l *Logger) LogSkip(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "SKIP: %s\n", path)
}

func (l *Logger) LogError(path string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.file, "ERROR: %s: %v\n", path, err)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
