package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetGitStatus runs git status --porcelain and returns a map of filename to status code.
func GetGitStatus(dirPath string) map[string]string {
	statusMap := make(map[string]string)
	
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dirPath
	var out bytes.Buffer
	cmd.Stdout = &out
	
	err := cmd.Run()
	if err != nil {
		return statusMap // Return empty map if error (not a git repo, or git not installed)
	}
	
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		// status is first two characters, then space, then filename
		status := strings.TrimSpace(line[:2])
		filename := strings.TrimSpace(line[3:])
		statusMap[filename] = status
	}
	
	return statusMap
}
