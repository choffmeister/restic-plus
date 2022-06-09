package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

type ExecCommandOpts struct {
	Name   string
	Args   []string
	Env    []string
	Logger *log.Logger
}

func ExecCommand(name string, args ...string) (string, int, error) {
	return ExecCommandWithOpts(ExecCommandOpts{
		Name: name,
		Args: args,
	})
}

func ExecCommandWithOpts(opts ExecCommandOpts) (string, int, error) {
	logger := opts.Logger
	if logger == nil {
		logger = LogDebug
	}

	LogDebug.Printf("Executing command %s %v\n", opts.Name, opts.Args)
	cmd := exec.Command(opts.Name, opts.Args...)

	var outputBuffer bytes.Buffer
	logWriter := NewLogWriter(logger)
	writer := io.MultiWriter(&outputBuffer, logWriter)
	cmd.Stdout = writer
	cmd.Stderr = writer
	cmd.Env = append(os.Environ(), opts.Env...)

	err := cmd.Run()
	outputStr := outputBuffer.String()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			return outputStr, 0, err
		}
		return outputStr, exitError.ExitCode(), fmt.Errorf("%w\n%s", exitError, outputStr)
	}
	return outputStr, 0, nil
}

func ExecCommandRetry(name string, args ...string) (string, int, error) {
	return ExecCommandRetryWithOpts(ExecCommandOpts{
		Name: name,
		Args: args,
	})
}

func ExecCommandRetryWithOpts(opts ExecCommandOpts) (string, int, error) {
	maxAttempts := 10
	delay := 1000 * time.Millisecond
	lastOutput := ""
	lastCode := 0
	lastErr := (error)(nil)
	attempt := 0
	for attempt < maxAttempts {
		if output, code, err := ExecCommandWithOpts(opts); err == nil {
			return output, code, err
		} else {
			lastOutput = output
			lastCode = code
			lastErr = err
		}
		attempt = attempt + 1
		time.Sleep(delay)
	}
	return lastOutput, lastCode, lastErr
}
