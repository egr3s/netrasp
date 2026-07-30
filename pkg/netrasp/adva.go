package netrasp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// ADVA is the Netrasp driver for Oscilloquartz OSA5401 CLI.
type adva struct {
	Connection connection
}

// Close connection to device.
func (a adva) Close(ctx context.Context) error {
	a.Connection.Close(ctx)

	return nil
}

// Configure device.
func (a adva) Configure(ctx context.Context, commands []string) (ConfigResult, error) {
	var result ConfigResult

	return result, fmt.Errorf("not implemented")
}

// Dial opens a connection to a device.
func (a adva) Dial(ctx context.Context) error {
	return establishConnection(ctx, a, a.Connection, a.basePrompt(), []string{})
}

// Enable elevates privileges.
func (a adva) Enable(ctx context.Context) error {
	return nil
}

// Run executes a command on a device.
func (a adva) Run(ctx context.Context, command string) (string, error) {
	output, err := a.RunUntil(ctx, command, a.basePrompt())
	if err != nil {
		return "", err
	}

	output = strings.ReplaceAll(output, "\r\n", "\n")
	lines := strings.Split(output, "\n")

	var sb strings.Builder

	for _, l := range lines {
		if l == "" {
			continue
		} else {
			sb.WriteString(l)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

// RunUntil executes a command and reads until the provided prompt.
func (a adva) RunUntil(ctx context.Context, command string, prompt *regexp.Regexp) (string, error) {
	err := a.Connection.Send(ctx, command)
	if err != nil {
		return "", fmt.Errorf("unable to send command to device: %w", err)
	}

	reader := a.Connection.Recv(ctx)
	output, err := readUntilPrompt(ctx, reader, prompt)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (a adva) basePrompt() *regexp.Regexp {
	return regexp.MustCompile(`(ADVA|ADTRAN)\S*-->`)
}
