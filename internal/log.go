package internal

import (
	"bytes"
	"log"
	"os"
	"sync"
)

var (
	LogDebug  = log.New(os.Stdout, "DEBUG:  ", 0)
	LogInfo   = log.New(os.Stdout, "INFO:   ", 0)
	LogWarn   = log.New(os.Stdout, "WARN:   ", 0)
	LogError  = log.New(os.Stdout, "ERROR:  ", log.Lshortfile)
	LogRestic = log.New(os.Stdout, "RESTIC: ", 0)
)

type LogWriter struct {
	mu     sync.Mutex
	logger *log.Logger
	buffer []byte
}

func NewLogWriter(logger *log.Logger) *LogWriter {
	return &LogWriter{logger: logger}
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	newline := byte('\n')
	breakline := []byte("\r")
	empty := []byte("")

	l.mu.Lock()
	defer l.mu.Unlock()

	l.buffer = append(l.buffer, bytes.ReplaceAll(p, breakline, empty)...)

	newlineIndex := bytes.IndexByte(l.buffer, newline)
	for newlineIndex >= 0 {
		line := l.buffer[0:newlineIndex]
		l.logger.Printf("%s\n", string(line))
		l.buffer = l.buffer[newlineIndex+1:]
		newlineIndex = bytes.IndexByte(l.buffer, newline)
	}

	return len(p), nil
}
