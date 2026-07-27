package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeDirectory(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "small.txt"), []byte("abc"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "large.txt"), []byte("12345678901234567890"), 0644)

	res, err := AnalyzeDirectory(tempDir)
	if err != nil {
		t.Fatalf("AnalyzeDirectory failed: %v", err)
	}

	if len(res.Items) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(res.Items))
	}

	if res.Items[0].Name != "large.txt" {
		t.Errorf("Expected largest item first, got %s", res.Items[0].Name)
	}

	if res.TotalSize != 23 {
		t.Errorf("Expected total size 23, got %d", res.TotalSize)
	}
}
