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
			fmt.Println("Matt Black Terminal File Manager v1.0.0")
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
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running Matt file manager: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Matt - Matt Black Terminal File Manager

Usage:
  matt [directory]
  matt [flags]

Flags:
  -h, --help      Show help and usage options
  -v, --version   Show version information

Keyboard Controls:
  Up / Down (k/j)   Navigate entries in focused pane
  Right / Enter (l) Open folder / expand directory
  Left (h)          Go up to parent directory
  Alt+D             Toggle Disk Space Analyzer view (ncdu style)
  /                 Fuzzy search and filter directory list
  Left Mouse Click  Focus pane & select item directly
  Tab / Shift+Tab   Cycle active focus between 3 upper panes & bottom terminal
  :                 Directly focus bottom terminal command runner / alias extensions
  .                 Toggle hidden files
  r                 Refresh directory view
  d                 Delete selected file/folder (triggers permission prompt)
  q / Ctrl+C        Quit Matt

Configuration & Extensions:
  Config file: ~/.config/matt/config.json
  Supports custom themes and command aliases (e.g. ll, findbig, count).`)
}
