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
			remaining = " " + strings.TrimSpace(cmdStr[len(head):])
		}
		cmdStr = strings.TrimSpace(aliasCmd + remaining)
		parts = strings.Fields(cmdStr)
		head = parts[0]
	}

	// Special handling for cd
	if head == "cd" {
		targetDir := config.GetHomeDir()
		if len(parts) > 1 {
			rawPath := strings.TrimSpace(cmdStr[2:])
			// Strip surrounding quotes if present
			if (strings.HasPrefix(rawPath, "\"") && strings.HasSuffix(rawPath, "\"")) ||
				(strings.HasPrefix(rawPath, "'") && strings.HasSuffix(rawPath, "'")) {
				rawPath = rawPath[1 : len(rawPath)-1]
			}

			expanded := config.ExpandPath(rawPath)
			targetDir = resolveDirectoryPath(expanded, currentDir)
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

	// Execute general shell command using user's shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", cmdStr)
	cmd.Dir = currentDir
	cmd.Env = os.Environ()

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

// resolveDirectoryPath intelligently finds target directory.
// Handles relative paths, missing leading slashes (e.g. "home/cachy"), and ~ home paths.
func resolveDirectoryPath(path string, currentDir string) string {
	if path == "" {
		return config.GetHomeDir()
	}

	// If absolute path and exists, return directly
	if filepath.IsAbs(path) {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return path
		}
	}

	// 1. Try path relative to currentDir
	relPath := filepath.Join(currentDir, path)
	if stat, err := os.Stat(relPath); err == nil && stat.IsDir() {
		return relPath
	}

	// 2. Try prepending "/" if user typed "home/cachy" without leading slash
	absWithSlash := "/" + strings.TrimPrefix(path, "/")
	if stat, err := os.Stat(absWithSlash); err == nil && stat.IsDir() {
		return absWithSlash
	}

	// 3. Fall back to relative path so standard OS error is returned
	return relPath
}
