package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ItemUsage holds size statistics for an item in the directory.
type ItemUsage struct {
	Name       string
	Path       string
	IsDir      bool
	SizeBytes  int64
	Percentage float64
	Bar        string
}

// DiskUsageResult contains analyzed disk breakdown.
type DiskUsageResult struct {
	TotalSize int64
	Items     []ItemUsage
	DirPath   string
}

// AnalyzeDirectory scans directory contents and computes relative sizes.
func AnalyzeDirectory(dirPath string) (DiskUsageResult, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return DiskUsageResult{}, err
	}

	var items []ItemUsage
	var totalSize int64

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		var size int64
		if entry.IsDir() {
			size = calculateDirSize(fullPath)
		} else {
			size = info.Size()
		}

		totalSize += size
		items = append(items, ItemUsage{
			Name:      entry.Name(),
			Path:      fullPath,
			IsDir:     entry.IsDir(),
			SizeBytes: size,
		})
	}

	// Sort largest to smallest
	sort.Slice(items, func(i, j int) bool {
		return items[i].SizeBytes > items[j].SizeBytes
	})

	// Calculate percentage and visual bar
	for i := range items {
		pct := float64(0)
		if totalSize > 0 {
			pct = (float64(items[i].SizeBytes) / float64(totalSize)) * 100
		}
		items[i].Percentage = pct
		items[i].Bar = generateBar(pct, 12)
	}

	return DiskUsageResult{
		TotalSize: totalSize,
		Items:     items,
		DirPath:   dirPath,
	}, nil
}

func calculateDirSize(dirPath string) int64 {
	var total int64
	_ = filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total
}

func generateBar(percentage float64, barLength int) string {
	filledLen := int((percentage / 100.0) * float64(barLength))
	if filledLen < 0 {
		filledLen = 0
	}
	if filledLen > barLength {
		filledLen = barLength
	}
	emptyLen := barLength - filledLen
	return strings.Repeat("█", filledLen) + strings.Repeat("░", emptyLen)
}

// FormatBytes converts raw byte count to readable string.
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
