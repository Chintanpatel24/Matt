package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Theme holds color tokens for the Matt Black aesthetic.
type Theme struct {
	Bg           string `json:"bg"`
	BgSurface    string `json:"bg_surface"`
	BgZebra      string `json:"bg_zebra"`
	Border       string `json:"border"`
	BorderActive string `json:"border_active"`
	TextPrimary  string `json:"text_primary"`
	TextMuted    string `json:"text_muted"`
	Directory    string `json:"directory"`
	Executable   string `json:"executable"`
	Selection    string `json:"selection"`
	Accent       string `json:"accent"`
	Error        string `json:"error"`
	Warning      string `json:"warning"`
	Success      string `json:"success"`
}

// Config represents application settings and custom command extension aliases.
type Config struct {
	ShowHidden bool              `json:"show_hidden"`
	Aliases    map[string]string `json:"aliases"`
	Theme      Theme             `json:"theme"`
	Bookmarks  []string          `json:"bookmarks"`
}

// DefaultConfig returns the standard Matt Black configuration with Lite Grey accents.
func DefaultConfig() Config {
	return Config{
		ShowHidden: false,
		Aliases: map[string]string{
			"ll":       "ls -la",
			"findbig":  "find . -type f -size +10M",
			"count":    "find . -type f | wc -l",
			"sysinfo":  "uname -a && uptime",
		},
		Theme: Theme{
			Bg:           "#09090b",
			BgSurface:    "#18181b",
			BgZebra:      "#0f0f12",
			Border:       "#27272a",
			BorderActive: "#71717a",
			TextPrimary:  "#f8fafc",
			TextMuted:    "#a1a1aa",
			Directory:    "#38bdf8",
			Executable:   "#34d399",
			Selection:    "#3f3f46",
			Accent:       "#a1a1aa",
			Error:        "#f87171",
			Warning:      "#fbbf24",
			Success:      "#4ade80",
		},
	}
}

// LoadConfig loads user config from ~/.config/matt/config.json and migrates any legacy yellow colors.
func LoadConfig() Config {
	cfg := DefaultConfig()
	configDir := filepath.Join(GetHomeDir(), ".config", "matt")
	configPath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		_ = os.MkdirAll(configDir, 0755)
		saveConfig(configPath, cfg)
		return cfg
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(data, &cfg)

	// Migrate old colors to new high-contrast gray theme
	modified := false
	if strings.EqualFold(cfg.Theme.BorderActive, "#eab308") || strings.EqualFold(cfg.Theme.BorderActive, "#e4e4e7") || cfg.Theme.BorderActive == "" {
		cfg.Theme.BorderActive = "#71717a"
		modified = true
	}
	if strings.EqualFold(cfg.Theme.Accent, "#eab308") || strings.EqualFold(cfg.Theme.Accent, "#e4e4e7") || cfg.Theme.Accent == "" {
		cfg.Theme.Accent = "#a1a1aa"
		modified = true
	}
	if strings.EqualFold(cfg.Theme.Selection, "#27272a") || cfg.Theme.Selection == "" {
		cfg.Theme.Selection = "#3f3f46"
		modified = true
	}
	if cfg.Theme.BgZebra == "" {
		cfg.Theme.BgZebra = "#0f0f12"
		modified = true
	}
	if cfg.Theme.Warning == "" {
		cfg.Theme.Warning = "#fbbf24"
		modified = true
	}

	if modified {
		saveConfig(configPath, cfg)
	}

	return cfg
}

func saveConfig(path string, cfg Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err == nil {
		_ = os.WriteFile(path, data, 0644)
	}
}

// GetHomeDir returns the user's home directory path.
func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return home
}

// ExpandPath expands leading ~ to home directory.
func ExpandPath(path string) string {
	if path == "~" {
		return GetHomeDir()
	}
	if len(path) > 1 && path[:2] == "~/" {
		return filepath.Join(GetHomeDir(), path[2:])
	}
	return path
}

// LoadBookmarks loads bookmarks from ~/.config/matt/bookmarks.json
func LoadBookmarks() []string {
	configDir := filepath.Join(GetHomeDir(), ".config", "matt")
	bookmarksPath := filepath.Join(configDir, "bookmarks.json")

	data, err := os.ReadFile(bookmarksPath)
	if err != nil {
		return []string{}
	}

	var bookmarks []string
	_ = json.Unmarshal(data, &bookmarks)
	if bookmarks == nil {
		bookmarks = []string{}
	}
	return bookmarks
}

// SaveBookmarks saves bookmarks to ~/.config/matt/bookmarks.json
func SaveBookmarks(bookmarks []string) {
	configDir := filepath.Join(GetHomeDir(), ".config", "matt")
	bookmarksPath := filepath.Join(configDir, "bookmarks.json")

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err == nil {
		_ = os.MkdirAll(configDir, 0755)
		_ = os.WriteFile(bookmarksPath, data, 0644)
	}
}

// AddBookmark adds a path to bookmarks if not duplicate
func AddBookmark(bookmarks []string, path string) []string {
	for _, b := range bookmarks {
		if b == path {
			return bookmarks
		}
	}
	return append(bookmarks, path)
}

// RemoveBookmark removes a bookmark by path
func RemoveBookmark(bookmarks []string, path string) []string {
	var result []string
	for _, b := range bookmarks {
		if b != path {
			result = append(result, b)
		}
	}
	if result == nil {
		result = []string{}
	}
	return result
}

// LoadHistory loads history from ~/.config/matt/history.json
func LoadHistory() []string {
	configDir := filepath.Join(GetHomeDir(), ".config", "matt")
	historyPath := filepath.Join(configDir, "history.json")

	data, err := os.ReadFile(historyPath)
	if err != nil {
		return []string{}
	}

	var history []string
	_ = json.Unmarshal(data, &history)
	if history == nil {
		history = []string{}
	}
	return history
}

// SaveHistory saves history to ~/.config/matt/history.json (max 100 entries)
func SaveHistory(history []string) {
	configDir := filepath.Join(GetHomeDir(), ".config", "matt")
	historyPath := filepath.Join(configDir, "history.json")

	if len(history) > 100 {
		history = history[:100]
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err == nil {
		_ = os.MkdirAll(configDir, 0755)
		_ = os.WriteFile(historyPath, data, 0644)
	}
}

// AddHistory adds cmd to front of history, deduplicates, and caps at 100
func AddHistory(history []string, cmd string) []string {
	var result []string
	result = append(result, cmd)
	for _, h := range history {
		if h != cmd {
			result = append(result, h)
		}
		if len(result) >= 100 {
			break
		}
	}
	return result
}
