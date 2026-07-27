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
	Border       string `json:"border"`
	BorderActive string `json:"border_active"`
	TextPrimary  string `json:"text_primary"`
	TextMuted    string `json:"text_muted"`
	Directory    string `json:"directory"`
	Executable   string `json:"executable"`
	Selection    string `json:"selection"`
	Accent       string `json:"accent"`
	Error        string `json:"error"`
	Success      string `json:"success"`
}

// Config represents application settings and custom command extension aliases.
type Config struct {
	ShowHidden bool              `json:"show_hidden"`
	Aliases    map[string]string `json:"aliases"`
	Theme      Theme             `json:"theme"`
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
			Bg:           "#09090b", // Deep Matte Black
			BgSurface:    "#18181b", // Dark Surface
			Border:       "#27272a", // Charcoal Border
			BorderActive: "#e4e4e7", // Lite Grey Active Accent
			TextPrimary:  "#f8fafc", // Pure White Primary Text
			TextMuted:    "#a1a1aa", // Soft Grey Muted Text
			Directory:    "#38bdf8", // Soft Cyan Directories
			Executable:   "#34d399", // Emerald Executables
			Selection:    "#27272a", // Active Selection Highlight
			Accent:       "#e4e4e7", // Lite Grey Accent
			Error:        "#f87171", // Red Error
			Success:      "#4ade80", // Green Success
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

	// Migrate old yellow colors (#eab308) to Lite Grey (#e4e4e7)
	modified := false
	if strings.EqualFold(cfg.Theme.BorderActive, "#eab308") {
		cfg.Theme.BorderActive = "#e4e4e7"
		modified = true
	}
	if strings.EqualFold(cfg.Theme.Accent, "#eab308") {
		cfg.Theme.Accent = "#e4e4e7"
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
