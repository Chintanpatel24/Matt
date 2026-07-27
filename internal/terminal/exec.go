package terminal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Chintanpatel24/Matt/internal/config"
)

// CommandResult stores output of command execution in Matt's bottom terminal.
type CommandResult struct {
	Output      string
	Err         error
	NewDir      string // If non-empty, tells Matt app to change current directory
	IsCdCommand bool
}

// ExecuteCommand runs a user command within the context of currentDir, resolving config aliases.
func ExecuteCommand(input string, currentDir string, cfg config.Config) CommandResult {
	cmdStr := strings.TrimSpace(input)
	if cmdStr == "" {
		return CommandResult{}
	}

	parts := strings.Fields(cmdStr)
	head := parts[0]

	// Check if alias exists in configuration
	if aliasCmd, exists := cfg.Aliases[head]; exists {
		remaining := ""
		if len(parts) > 1 {
			remaining = " " + strings.Join(parts[1:], " ")
		}
		cmdStr = aliasCmd + remaining
		parts = strings.Fields(cmdStr)
		head = parts[0]
	}

	// Special handling for cd
	if head == "cd" {
		targetDir := config.GetHomeDir()
		if len(parts) > 1 {
			targetDir = config.ExpandPath(parts[1])
			if !filepath.IsAbs(targetDir) {
				targetDir = filepath.Join(currentDir, targetDir)
			}
		}
		targetDir = filepath.Clean(targetDir)

		stat, err := os.Stat(targetDir)
		if err != nil {
			return CommandResult{
				Output: fmt.Sprintf("cd: %s: No such file or directory", targetDir),
				Err:    err,
			}
		}
		if !stat.IsDir() {
			return CommandResult{
				Output: fmt.Sprintf("cd: %s: Not a directory", targetDir),
				Err:    fmt.Errorf("not a directory"),
			}
		}

		return CommandResult{
			Output:      fmt.Sprintf("Changed directory to %s", targetDir),
			NewDir:      targetDir,
			IsCdCommand: true,
		}
	}

	// Execute general shell command using bash/sh
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", cmdStr)
	cmd.Dir = currentDir

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	outputStr := strings.TrimRight(outBuf.String(), "\r\n")

	if err != nil && outputStr == "" {
		outputStr = fmt.Sprintf("Command failed: %v", err)
	}

	return CommandResult{
		Output: outputStr,
		Err:    err,
	}
}
