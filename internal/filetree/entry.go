package filetree

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// FileEntry represents a file or directory node in Matt file manager.
type FileEntry struct {
	Name          string
	Path          string
	IsDir         bool
	IsSymlink     bool
	Size          int64
	Mode          os.FileMode
	ModTime       time.Time
	Extension     string
	Icon          string
	IsHidden      bool
	Permissions   string
	Owner         string
	Group         string
	SymlinkTarget string
	GitStatus     string // e.g. "M", "A", "?", or ""
}

// NewFileEntry creates a FileEntry struct from os.FileInfo.
func NewFileEntry(dirPath string, info os.FileInfo) FileEntry {
	name := info.Name()
	fullPath := filepath.Join(dirPath, name)
	isDir := info.IsDir()
	isSymlink := info.Mode()&os.ModeSymlink != 0
	ext := strings.ToLower(filepath.Ext(name))
	isHidden := strings.HasPrefix(name, ".")

	owner, group := getOwnerGroup(info)

	entry := FileEntry{
		Name:        name,
		Path:        fullPath,
		IsDir:       isDir,
		IsSymlink:   isSymlink,
		Size:        info.Size(),
		Mode:        info.Mode(),
		ModTime:     info.ModTime(),
		Extension:   ext,
		Icon:        GetIcon(name, isDir, ext),
		IsHidden:    isHidden,
		Permissions: info.Mode().String(),
		Owner:       owner,
		Group:       group,
	}

	// Resolve symlink target
	if isSymlink {
		entry.Icon = "🔗"
		target, err := os.Readlink(fullPath)
		if err == nil {
			entry.SymlinkTarget = target
		}
	}

	return entry
}

func getOwnerGroup(info os.FileInfo) (string, string) {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		u := strconv.FormatUint(uint64(stat.Uid), 10)
		g := strconv.FormatUint(uint64(stat.Gid), 10)

		usr, err := user.LookupId(u)
		if err == nil {
			u = usr.Username
		}
		grp, err := user.LookupGroupId(g)
		if err == nil {
			g = grp.Name
		}
		return u, g
	}
	return "user", "group"
}

// FormatSize returns human-readable file size string.
func (f FileEntry) FormatSize() string {
	if f.IsDir {
		return "<DIR>"
	}
	const unit = 1024
	if f.Size < unit {
		return fmt.Sprintf("%d B", f.Size)
	}
	div, exp := int64(unit), 0
	for n := f.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(f.Size)/float64(div), "KMGTPE"[exp])
}

// FormatModTime returns formatted modification date/time.
func (f FileEntry) FormatModTime() string {
	return f.ModTime.Format("2006-01-02 15:04")
}

// FormatSymlink returns a display string showing the symlink target.
func (f FileEntry) FormatSymlink() string {
	if f.IsSymlink && f.SymlinkTarget != "" {
		return fmt.Sprintf("→ %s", f.SymlinkTarget)
	}
	return ""
}

// GetIcon returns an appropriate unicode icon for file types.
func GetIcon(name string, isDir bool, ext string) string {
	if isDir {
		if name == ".." || name == "." {
			return "↩ "
		}
		return "📁"
	}

	switch ext {
	case ".go":
		return "🐹"
	case ".rs":
		return "🦀"
	case ".py":
		return "🐍"
	case ".js", ".ts", ".jsx", ".tsx":
		return "🟨"
	case ".html", ".htm":
		return "🌐"
	case ".css", ".scss", ".less":
		return "🎨"
	case ".json", ".yaml", ".yml", ".toml", ".xml":
		return "⚙️"
	case ".md", ".txt", ".doc", ".docx":
		return "📝"
	case ".sh", ".bash", ".zsh":
		return "📜"
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return "🖼️"
	case ".mp3", ".wav", ".flac", ".ogg":
		return "🎵"
	case ".mp4", ".mkv", ".avi", ".mov":
		return "🎬"
	case ".zip", ".tar", ".gz", ".7z", ".rar":
		return "📦"
	case ".pdf":
		return "📕"
	case ".db", ".sqlite", ".sql":
		return "🗄️"
	default:
		return "📄"
	}
}
