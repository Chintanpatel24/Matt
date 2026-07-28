package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Chintanpatel24/Matt/internal/app"
	"github.com/Chintanpatel24/Matt/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	targetDir, err := os.Getwd()
	if err != nil {
		targetDir = "."
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-h" || arg == "--help" {
			printUsage()
			os.Exit(0)
		}
		if arg == "-v" || arg == "--version" {
			fmt.Println("Matt Black Terminal File Manager v2.0.0")
			os.Exit(0)
		}
		expanded := config.ExpandPath(arg)
		absPath, err := filepath.Abs(expanded)
		if err == nil {
			stat, err := os.Stat(absPath)
			if err == nil && stat.IsDir() {
				targetDir = absPath
			}
		}
	}

	initialModel := app.NewAppModel(cfg, targetDir)

	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running Matt file manager: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Matt - Matt Black Terminal File Manager v2.0.0

Usage:
  matt [directory]
  matt [flags]

Flags:
  -h, --help      Show help and usage options
  -v, --version   Show version information

Keyboard Controls:
  Up / Down (k/j)   Navigate entries in focused pane
  Right / Enter (l)  Open folder / expand directory
  Left (h)           Go up to parent directory
  Alt+D              Toggle Disk Space Analyzer view (async)
  /                  Fuzzy search and filter directory list
  Tab / Shift+Tab    Cycle active focus between panes & terminal
  :                  Focus bottom terminal command runner
  .                  Toggle hidden files
  r                  Refresh directory view
  d                  Delete selected file/folder (permission prompt)
  n                  Create new file or folder
  c                  Copy selected file/folder
  p                  Paste copied item into current directory
  m                  Move/rename selected file/folder
  b                  Open bookmarks list
  B                  Bookmark current directory
  g                  Jump to first item
  G                  Jump to last item
  q / Ctrl+C         Quit Matt

Configuration & Extensions:
  Config file: ~/.config/matt/config.json
  Supports custom themes and command aliases (e.g. ll, findbig, count).`)
}
