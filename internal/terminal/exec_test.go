package terminal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Chintanpatel24/Matt/internal/config"
)

func TestExecuteCommand_Cd(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "target")
	_ = os.Mkdir(subDir, 0755)

	cfg := config.DefaultConfig()
	res := ExecuteCommand("cd target", tempDir, cfg)
	if !res.IsCdCommand {
		t.Errorf("Expected IsCdCommand to be true")
	}
	if res.NewDir != subDir {
		t.Errorf("Expected NewDir to be %s, got %s", subDir, res.NewDir)
	}
}

func TestExecuteCommand_Alias(t *testing.T) {
	tempDir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.Aliases["hello"] = "echo 'hello from alias'"

	res := ExecuteCommand("hello", tempDir, cfg)
	if res.Err != nil {
		t.Fatalf("Unexpected command error: %v", res.Err)
	}
	if res.Output != "hello from alias" {
		t.Errorf("Expected 'hello from alias', got '%s'", res.Output)
	}
}

func TestExecuteCommand_Dot(t *testing.T) {
	tempDir := t.TempDir()
	cfg := config.DefaultConfig()

	res := ExecuteCommand(".", tempDir, cfg)
	if res.Output == "" {
		t.Errorf("Expected output from '.' command, got empty string")
	}
}

