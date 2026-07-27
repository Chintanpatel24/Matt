package filetree

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	tempDir := t.TempDir()

	// Create test structure
	_ = os.Mkdir(filepath.Join(tempDir, "subfolder"), 0755)
	_ = os.WriteFile(filepath.Join(tempDir, "hello.txt"), []byte("hello world"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, ".hidden"), []byte("secret"), 0644)

	// Test without hidden files
	entries, err := ReadDir(tempDir, false)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	foundSubfolder := false
	foundHello := false
	foundHidden := false

	for _, e := range entries {
		if e.Name == "subfolder" && e.IsDir {
			foundSubfolder = true
		}
		if e.Name == "hello.txt" && !e.IsDir {
			foundHello = true
		}
		if e.Name == ".hidden" {
			foundHidden = true
		}
	}

	if !foundSubfolder {
		t.Errorf("Expected subfolder in ReadDir output")
	}
	if !foundHello {
		t.Errorf("Expected hello.txt in ReadDir output")
	}
	if foundHidden {
		t.Errorf("Did not expect hidden file when showHidden is false")
	}

	// Test with hidden files
	entriesHidden, err := ReadDir(tempDir, true)
	if err != nil {
		t.Fatalf("ReadDir with hidden failed: %v", err)
	}

	foundHidden = false
	for _, e := range entriesHidden {
		if e.Name == ".hidden" {
			foundHidden = true
		}
	}
	if !foundHidden {
		t.Errorf("Expected hidden file when showHidden is true")
	}
}

func TestFileEntryFormatting(t *testing.T) {
	fe := FileEntry{
		Name: "test.png",
		Size: 2048,
	}
	if fe.FormatSize() != "2.0 KB" {
		t.Errorf("Expected '2.0 KB', got '%s'", fe.FormatSize())
	}

	icon := GetIcon("main.go", false, ".go")
	if icon != "🐹" {
		t.Errorf("Expected 🐹 icon for .go file, got %s", icon)
	}
}
