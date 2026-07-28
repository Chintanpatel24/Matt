package filetree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Chintanpatel24/Matt/internal/git"
)

// ReadDir reads directory content and returns sorted FileEntry slice.
func ReadDir(dirPath string, showHidden bool) ([]FileEntry, error) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		absPath = dirPath
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	gitStatuses := git.GetGitStatus(absPath)

	var results []FileEntry

	// Add parent directory option if not at system root
	if absPath != "/" {
		parentPath := filepath.Dir(absPath)
		parentInfo, err := os.Stat(parentPath)
		if err == nil {
			parentEntry := NewFileEntry(filepath.Dir(parentPath), parentInfo)
			parentEntry.Name = ".."
			parentEntry.Path = parentPath
			parentEntry.Icon = "↩ "
			results = append(results, parentEntry)
		}
	}

	var dirs []FileEntry
	var files []FileEntry

	for _, entry := range entries {
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fe := NewFileEntry(absPath, info)
		if status, exists := gitStatuses[fe.Name]; exists {
			fe.GitStatus = status
		}
		if fe.IsDir {
			dirs = append(dirs, fe)
		} else {
			files = append(files, fe)
		}
	}

	// Sort directories and files alphabetically case-insensitive
	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	results = append(results, dirs...)
	results = append(results, files...)

	return results, nil
}
