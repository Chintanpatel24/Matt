package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Chintanpatel24/Matt/internal/analyzer"
	"github.com/Chintanpatel24/Matt/internal/config"
	"github.com/Chintanpatel24/Matt/internal/filetree"
	"github.com/Chintanpatel24/Matt/internal/fuzzy"
	"github.com/Chintanpatel24/Matt/internal/preview"
	"github.com/Chintanpatel24/Matt/internal/terminal"
	"github.com/Chintanpatel24/Matt/internal/ui"
	"github.com/Chintanpatel24/Matt/internal/version"
)

// ActiveFocus specifies which pane currently has key focus.
type ActiveFocus int

const (
	PaneLeft ActiveFocus = iota
	PaneCenter
	PaneRight
	PaneTerminal
)

// AppState describes root application state.
type AppState int

const (
	StateStartup AppState = iota
	StateMain
	StateDialog
	StateAnalyzer
	StateAnalyzerLoading
	StateBookmarks
	StateInput
)

// InputAction describes what the text input prompt is for.
type InputAction int

const (
	InputNewFile InputAction = iota
	InputNewFolder
	InputRename
)

// UndoAction describes an action that can be reverted.
type UndoAction struct {
	Type   string // "rename", "create", "delete", "paste", "archive"
	Source string
	Dest   string
}

// AppModel is the primary Elm Architecture state for Matt.
type AppModel struct {
	State         AppState
	Config        config.Config
	Styles        ui.Styles
	Width         int
	Height        int
	CurrentDir    string
	AllEntries    []filetree.FileEntry
	LeftEntries   []filetree.FileEntry
	LeftCursor    int
	LeftScroll    int
	CenterEntries []filetree.FileEntry
	CenterCursor  int
	CenterScroll  int
	RightPreview  preview.PreviewResult
	RightScroll   int
	SelectedItem  filetree.FileEntry
	Focus         ActiveFocus

	// Inputs
	TermInput    textinput.Model
	FilterInput  textinput.Model
	PromptInput  textinput.Model
	IsFiltering  bool

	// Analyzer state
	AnalyzerResult analyzer.DiskUsageResult
	AnalyzerCursor int
	AnalyzerScroll int

	// Command history
	CmdHistory    []string
	HistoryIdx    int
	HistorySaved  string

	// Clipboard for copy/paste
	ClipboardPaths []string // Supports multi-selection copying/cutting
	ClipboardOp    string   // "copy" or "cut"

	// Multi-select state
	MultiSelect map[string]bool // Maps file path to true if selected

	// Undo / Redo stack
	UndoStack []UndoAction

	// Bookmarks
	Bookmarks      []string
	BookmarkCursor int
	BookmarkScroll int

	// Input prompt state
	InputAction InputAction

	LastCmdOut string
	StatusMsg  string

	// Sub-models
	Startup StartupModel
	Dialog  ui.PermissionDialog
}

// NewAppModel initializes Matt application state.
func NewAppModel(cfg config.Config, initialDir string) AppModel {
	absDir, err := filepath.Abs(initialDir)
	if err != nil {
		absDir = initialDir
	}

	ti := textinput.New()
	ti.Placeholder = "Type shell command or alias (e.g. cd, touch, findbig, ll)..."
	ti.Prompt = "matt $ "
	ti.CharLimit = 256
	ti.Width = 60

	fi := textinput.New()
	fi.Placeholder = "Fuzzy search files..."
	fi.Prompt = "🔍 / "
	fi.CharLimit = 64
	fi.Width = 24

	pi := textinput.New()
	pi.Placeholder = "Enter name..."
	pi.Prompt = "▌ "
	pi.CharLimit = 256
	pi.Width = 40

	styles := ui.NewStyles(cfg)
	startup := NewStartupModel(cfg, absDir)

	m := AppModel{
		State:          StateStartup,
		Config:         cfg,
		Styles:         styles,
		CurrentDir:     absDir,
		Focus:          PaneLeft,
		TermInput:      ti,
		FilterInput:    fi,
		PromptInput:    pi,
		IsFiltering:    false,
		Startup:        startup,
		CmdHistory:     config.LoadHistory(),
		HistoryIdx:     -1,
		Bookmarks:      config.LoadBookmarks(),
		MultiSelect:    make(map[string]bool),
		UndoStack:      []UndoAction{},
		ClipboardPaths: []string{},
		LastCmdOut:     "Ready. Press [Tab] focus • [/] filter • [Alt+D] analyzer • [:] cmd • [.] hidden",
		StatusMsg:      fmt.Sprintf("Matt %s — Matt Black Terminal File Manager", version.Version),
	}

	m.refreshLeftEntries()
	m.refreshCenterAndRight()

	return m
}

func (m *AppModel) refreshLeftEntries() {
	entries, err := filetree.ReadDir(m.CurrentDir, m.Config.ShowHidden)
	if err != nil {
		m.StatusMsg = fmt.Sprintf("Error reading directory: %v", err)
		return
	}
	m.AllEntries = entries

	query := strings.TrimSpace(m.FilterInput.Value())
	if query != "" {
		var filtered []filetree.FileEntry
		for _, e := range entries {
			if e.Name == ".." || fuzzy.Match(query, e.Name) {
				filtered = append(filtered, e)
			}
		}
		m.LeftEntries = filtered
	} else {
		m.LeftEntries = entries
	}

	if m.LeftCursor >= len(m.LeftEntries) {
		m.LeftCursor = max(0, len(m.LeftEntries)-1)
	}
}

func (m *AppModel) refreshCenterAndRight() {
	if len(m.LeftEntries) == 0 {
		m.CenterEntries = nil
		m.RightPreview = preview.PreviewResult{Content: "\n  [No files in view]"}
		m.SelectedItem = filetree.FileEntry{}
		return
	}

	if m.LeftCursor >= len(m.LeftEntries) {
		m.LeftCursor = 0
	}
	if m.LeftCursor < 0 {
		m.LeftCursor = 0
	}

	selected := m.LeftEntries[m.LeftCursor]

	if selected.IsDir {
		entries, err := filetree.ReadDir(selected.Path, m.Config.ShowHidden)
		if err == nil {
			m.CenterEntries = entries
		} else {
			m.CenterEntries = nil
		}
	} else {
		m.CenterEntries = m.LeftEntries
	}

	if m.Focus == PaneLeft {
		m.CenterCursor = 0
		m.CenterScroll = 0
	} else {
		if m.CenterCursor >= len(m.CenterEntries) {
			m.CenterCursor = max(0, len(m.CenterEntries)-1)
		}
		if m.CenterCursor < 0 {
			m.CenterCursor = 0
		}
	}

	if m.Focus == PaneCenter && len(m.CenterEntries) > 0 {
		selected = m.CenterEntries[m.CenterCursor]
	}

	m.SelectedItem = selected

	previewWidth := max(30, (m.Width/3)-4)
	previewHeight := max(6, m.Height-18)
	m.RightPreview = preview.GeneratePreview(selected, previewWidth, previewHeight)
}

// countDirsFiles returns separate dir and file counts.
func countDirsFiles(entries []filetree.FileEntry) (int, int) {
	dirs, files := 0, 0
	for _, e := range entries {
		if e.Name == ".." {
			continue
		}
		if e.IsDir {
			dirs++
		} else {
			files++
		}
	}
	return dirs, files
}

// buildBreadcrumb generates a breadcrumb path display.
func (m AppModel) buildBreadcrumb() string {
	home := config.GetHomeDir()
	path := m.CurrentDir

	// Shorten home directory to ~
	if strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}

	parts := strings.Split(path, "/")
	var rendered []string
	for i, part := range parts {
		if part == "" {
			if i == 0 {
				part = "/"
			} else {
				continue
			}
		}
		if i == len(parts)-1 {
			rendered = append(rendered, m.Styles.BreadcrumbActive.Render(part))
		} else {
			rendered = append(rendered, m.Styles.Breadcrumb.Render(part))
		}
	}
	return strings.Join(rendered, m.Styles.HeaderPathSep.Render(" ❯ "))
}

// getModeIndicator returns the styled mode badge.
func (m AppModel) getModeIndicator() string {
	if m.IsFiltering {
		return m.Styles.ModeFilter.Render(" FILTER ")
	}
	if m.Focus == PaneTerminal {
		return m.Styles.ModeCommand.Render(" COMMAND ")
	}
	return m.Styles.ModeNormal.Render(" NORMAL ")
}

func (m AppModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Startup.Width = msg.Width
		m.Startup.Height = msg.Height
		m.TermInput.Width = max(20, msg.Width-15)
		m.refreshCenterAndRight()
		return m, nil

	case analyzer.AnalyzerResultMsg:
		if msg.Err != nil {
			m.StatusMsg = fmt.Sprintf("Analyzer error: %v", msg.Err)
			m.State = StateMain
		} else {
			m.AnalyzerResult = msg.Result
			m.AnalyzerCursor = 0
			m.AnalyzerScroll = 0
			m.State = StateAnalyzer
		}
		return m, nil
	}

	// 1. Startup Mode
	if m.State == StateStartup {
		m.Startup, cmd = m.Startup.Update(msg)
		if m.Startup.IsDone {
			if m.Startup.Selected == ChoiceHomeDir {
				m.CurrentDir = config.GetHomeDir()
			} else {
				m.CurrentDir = m.Startup.CurrentPath
			}
			m.State = StateMain
			m.LeftCursor = 0
			m.LeftScroll = 0
			m.CenterCursor = 0
			m.CenterScroll = 0
			m.RightScroll = 0
			m.refreshLeftEntries()
			m.refreshCenterAndRight()
		}
		return m, cmd
	}

	// 2. Dialog Mode
	if m.State == StateDialog {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "left", "h":
				m.Dialog.ActiveIndex = 0
			case "right", "l":
				m.Dialog.ActiveIndex = 1
			case "enter":
				if m.Dialog.ActiveIndex == 0 && m.Dialog.OnConfirm != nil {
					m.Dialog.OnConfirm()
				}
				m.State = StateMain
				m.Dialog.IsOpen = false
			case "esc", "q":
				m.State = StateMain
				m.Dialog.IsOpen = false
			}
		}
		return m, nil
	}

	// 3. Input Prompt Mode (new file, new folder, rename)
	if m.State == StateInput {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.State = StateMain
				m.PromptInput.Blur()
				m.PromptInput.SetValue("")
				return m, nil
			case "enter":
				name := strings.TrimSpace(m.PromptInput.Value())
				if name != "" {
					switch m.InputAction {
					case InputNewFile:
						newPath := filepath.Join(m.CurrentDir, name)
						err := os.WriteFile(newPath, []byte{}, 0644)
						if err != nil {
							m.LastCmdOut = m.Styles.ErrorText.Render(fmt.Sprintf("✗ Error creating file: %v", err))
						} else {
							m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Created file: %s", name))
						}
					case InputNewFolder:
						newPath := filepath.Join(m.CurrentDir, name)
						err := os.MkdirAll(newPath, 0755)
						if err != nil {
							m.LastCmdOut = m.Styles.ErrorText.Render(fmt.Sprintf("✗ Error creating folder: %v", err))
						} else {
							m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Created folder: %s", name))
						}
					case InputRename:
						if len(m.LeftEntries) > 0 {
							target := m.LeftEntries[m.LeftCursor]
							newPath := filepath.Join(m.CurrentDir, name)
							err := os.Rename(target.Path, newPath)
							if err != nil {
								m.LastCmdOut = m.Styles.ErrorText.Render(fmt.Sprintf("✗ Error renaming: %v", err))
							} else {
								m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Renamed '%s' → '%s'", target.Name, name))
								m.UndoStack = append(m.UndoStack, UndoAction{Type: "rename", Source: newPath, Dest: target.Path})
							}
						}
					}
					m.refreshLeftEntries()
					m.refreshCenterAndRight()
				}
				m.State = StateMain
				m.PromptInput.Blur()
				m.PromptInput.SetValue("")
				return m, nil
			default:
				m.PromptInput, cmd = m.PromptInput.Update(msg)
				return m, cmd
			}
		}
		return m, nil
	}

	// 4. Bookmarks Mode
	if m.State == StateBookmarks {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc", "b", "q":
				m.State = StateMain
				return m, nil
			case "up", "k":
				if m.BookmarkCursor > 0 {
					m.BookmarkCursor--
				}
			case "down", "j":
				if m.BookmarkCursor < len(m.Bookmarks)-1 {
					m.BookmarkCursor++
				}
			case "enter":
				if len(m.Bookmarks) > 0 && m.BookmarkCursor < len(m.Bookmarks) {
					target := m.Bookmarks[m.BookmarkCursor]
					stat, err := os.Stat(target)
					if err == nil && stat.IsDir() {
						m.CurrentDir = target
						m.LeftCursor = 0
						m.LeftScroll = 0
						m.CenterCursor = 0
						m.CenterScroll = 0
						m.RightScroll = 0
						m.FilterInput.SetValue("")
						m.refreshLeftEntries()
						m.refreshCenterAndRight()
						m.StatusMsg = fmt.Sprintf("Jumped to: %s", target)
					} else {
						m.LastCmdOut = m.Styles.ErrorText.Render(fmt.Sprintf("✗ Bookmark path invalid: %s", target))
					}
					m.State = StateMain
				}
			case "d", "x":
				if len(m.Bookmarks) > 0 && m.BookmarkCursor < len(m.Bookmarks) {
					removed := m.Bookmarks[m.BookmarkCursor]
					m.Bookmarks = config.RemoveBookmark(m.Bookmarks, removed)
					config.SaveBookmarks(m.Bookmarks)
					if m.BookmarkCursor >= len(m.Bookmarks) {
						m.BookmarkCursor = max(0, len(m.Bookmarks)-1)
					}
					m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Removed bookmark: %s", removed))
				}
			}
		}
		return m, nil
	}

	// 5. Disk Space Analyzer View Mode
	if m.State == StateAnalyzer {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc", "alt+d", "q":
				m.State = StateMain
				return m, nil
			case "up", "k":
				if m.AnalyzerCursor > 0 {
					m.AnalyzerCursor--
				}
			case "down", "j":
				if m.AnalyzerCursor < len(m.AnalyzerResult.Items)-1 {
					m.AnalyzerCursor++
				}
			case "s":
				m.AnalyzerResult.ToggleSort()
				m.AnalyzerCursor = 0
				m.AnalyzerScroll = 0
			case "enter":
				if len(m.AnalyzerResult.Items) > 0 {
					selected := m.AnalyzerResult.Items[m.AnalyzerCursor]
					if selected.IsDir {
						m.CurrentDir = selected.Path
						m.State = StateMain
						m.LeftCursor = 0
						m.LeftScroll = 0
						m.CenterCursor = 0
						m.CenterScroll = 0
						m.RightScroll = 0
						m.refreshLeftEntries()
						m.refreshCenterAndRight()
					}
				}
			}
		}
		return m, nil
	}

	// 6. Analyzer Loading State
	if m.State == StateAnalyzerLoading {
		return m, nil
	}

	// 7. Main View Mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "alt+d" || msg.String() == "alt+D" {
			m.State = StateAnalyzerLoading
			m.StatusMsg = "⏳ Analyzing disk usage..."
			return m, analyzer.AnalyzeDirectoryAsync(m.CurrentDir)
		}

		if m.IsFiltering {
			switch msg.String() {
			case "esc":
				m.IsFiltering = false
				m.FilterInput.Blur()
				m.FilterInput.SetValue("")
				m.LeftCursor = 0
				m.LeftScroll = 0
				m.refreshLeftEntries()
				m.refreshCenterAndRight()
				return m, nil

			case "enter":
				m.IsFiltering = false
				m.FilterInput.Blur()
				return m, nil
			default:
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				m.LeftCursor = 0
				m.LeftScroll = 0
				m.refreshLeftEntries()
				m.refreshCenterAndRight()
				return m, cmd
			}
		}

		if m.Focus != PaneTerminal {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "/":
				m.IsFiltering = true
				m.FilterInput.Focus()
				return m, nil

			case "tab":
				m.Focus = (m.Focus + 1) % 4
				if m.Focus == PaneTerminal {
					m.TermInput.Focus()
				} else {
					m.TermInput.Blur()
				}
				m.refreshCenterAndRight()
				return m, nil

			case "shift+tab":
				m.Focus = (m.Focus + 3) % 4
				if m.Focus == PaneTerminal {
					m.TermInput.Focus()
				} else {
					m.TermInput.Blur()
				}
				m.refreshCenterAndRight()
				return m, nil

			case ":":
				m.Focus = PaneTerminal
				m.TermInput.Focus()
				return m, nil

			case ".":
				m.Config.ShowHidden = !m.Config.ShowHidden
				m.refreshLeftEntries()
				m.refreshCenterAndRight()
				if m.Config.ShowHidden {
					m.StatusMsg = "Hidden files: shown"
				} else {
					m.StatusMsg = "Hidden files: hidden"
				}
				return m, nil

			case "r":
				m.refreshLeftEntries()
				m.refreshCenterAndRight()
				m.StatusMsg = "✓ Refreshed view"
				return m, nil

			case " ":
				if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						if m.MultiSelect[target.Path] {
							delete(m.MultiSelect, target.Path)
						} else {
							m.MultiSelect[target.Path] = true
						}
						m.StatusMsg = fmt.Sprintf("Selected: %d items", len(m.MultiSelect))
					}
				}
				return m, nil

			case "d":
				var targets []string
				if len(m.MultiSelect) > 0 {
					for p := range m.MultiSelect {
						targets = append(targets, p)
					}
				} else if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						targets = append(targets, target.Path)
					}
				}

				if len(targets) > 0 {
					confirmMsg := fmt.Sprintf("Are you sure you want to delete %d items?", len(targets))
					if len(targets) == 1 {
						confirmMsg = fmt.Sprintf("Are you sure you want to delete '%s'?", filepath.Base(targets[0]))
					}
					m.Dialog = ui.NewPermissionDialog(
						"Delete Confirmation",
						confirmMsg,
						func() {
							deletedCount := 0
							for _, targetPath := range targets {
								err := os.RemoveAll(targetPath)
								if err == nil {
									deletedCount++
								}
							}
							m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Deleted %d items", deletedCount))
							m.MultiSelect = make(map[string]bool)
							m.refreshLeftEntries()
							m.refreshCenterAndRight()
						},
					)
					m.State = StateDialog
					return m, nil
				}

			case "n":
				m.InputAction = InputNewFile
				m.PromptInput.Placeholder = "Enter filename..."
				m.PromptInput.SetValue("")
				m.PromptInput.Focus()
				m.State = StateInput
				return m, nil

			case "N":
				m.InputAction = InputNewFolder
				m.PromptInput.Placeholder = "Enter folder name..."
				m.PromptInput.SetValue("")
				m.PromptInput.Focus()
				m.State = StateInput
				return m, nil

			case "m":
				if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						m.InputAction = InputRename
						m.PromptInput.Placeholder = fmt.Sprintf("Rename '%s' to...", target.Name)
						m.PromptInput.SetValue(target.Name)
						m.PromptInput.Focus()
						m.State = StateInput
						return m, nil
					}
				}

			case "c", "x":
				opName := "copy"
				if msg.String() == "x" {
					opName = "cut"
				}
				m.ClipboardOp = opName
				m.ClipboardPaths = []string{}
				
				if len(m.MultiSelect) > 0 {
					for path := range m.MultiSelect {
						m.ClipboardPaths = append(m.ClipboardPaths, path)
					}
				} else if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						m.ClipboardPaths = append(m.ClipboardPaths, target.Path)
					}
				}
				
				if len(m.ClipboardPaths) > 0 {
					m.StatusMsg = fmt.Sprintf("Clipboard: %d items marked for %s", len(m.ClipboardPaths), opName)
				} else {
					m.StatusMsg = "Nothing selected to copy/cut"
				}
				return m, nil

			case "p":
				if len(m.ClipboardPaths) > 0 {
					successCount := 0
					for _, src := range m.ClipboardPaths {
						srcName := filepath.Base(src)
						dest := filepath.Join(m.CurrentDir, srcName)
						
						var err error
						if m.ClipboardOp == "copy" {
							err = exec.Command("cp", "-r", src, dest).Run()
							if err == nil {
								m.UndoStack = append(m.UndoStack, UndoAction{Type: "paste", Dest: dest})
							}
						} else { // cut
							err = os.Rename(src, dest)
							if err == nil {
								m.UndoStack = append(m.UndoStack, UndoAction{Type: "rename", Source: dest, Dest: src})
							}
						}
						if err == nil {
							successCount++
						}
					}
					
					m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Pasted %d/%d items", successCount, len(m.ClipboardPaths)))
					if m.ClipboardOp == "cut" {
						m.ClipboardPaths = []string{}
						m.ClipboardOp = ""
					}
					m.MultiSelect = make(map[string]bool)
					m.refreshLeftEntries()
					m.refreshCenterAndRight()
				} else {
					m.StatusMsg = "Clipboard is empty. Copy/cut first."
				}
				return m, nil

			case "u":
				if len(m.UndoStack) > 0 {
					action := m.UndoStack[len(m.UndoStack)-1]
					m.UndoStack = m.UndoStack[:len(m.UndoStack)-1]
					
					var err error
					switch action.Type {
					case "rename":
						err = os.Rename(action.Source, action.Dest)
						if err == nil {
							m.StatusMsg = fmt.Sprintf("Undone: rename of %s", filepath.Base(action.Source))
						}
					case "paste":
						err = os.RemoveAll(action.Dest)
						if err == nil {
							m.StatusMsg = fmt.Sprintf("Undone: pasted %s", filepath.Base(action.Dest))
						}
					}
					if err != nil {
						m.StatusMsg = fmt.Sprintf("Undo failed: %v", err)
					}
					m.refreshLeftEntries()
					m.refreshCenterAndRight()
				} else {
					m.StatusMsg = "Nothing to undo"
				}
				return m, nil

			case "z":
				var targets []string
				if len(m.MultiSelect) > 0 {
					for p := range m.MultiSelect {
						targets = append(targets, p)
					}
				} else if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						targets = append(targets, target.Path)
					}
				}
				
				if len(targets) > 0 {
					zipPath := filepath.Join(m.CurrentDir, "archive.zip")
					args := []string{"-r", zipPath}
					for _, t := range targets {
						args = append(args, filepath.Base(t))
					}
					cmd := exec.Command("zip", args...)
					cmd.Dir = m.CurrentDir
					err := cmd.Run()
					if err != nil {
						m.LastCmdOut = m.Styles.ErrorText.Render(fmt.Sprintf("✗ Zip failed: %v (make sure zip is installed)", err))
					} else {
						m.LastCmdOut = m.Styles.SuccessText.Render("✓ Created archive.zip")
						m.MultiSelect = make(map[string]bool)
						m.refreshLeftEntries()
						m.refreshCenterAndRight()
					}
				} else {
					m.StatusMsg = "Nothing selected to compress"
				}
				return m, nil

			case "Z":
				if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						ext := strings.ToLower(filepath.Ext(target.Name))
						var cmd *exec.Cmd
						if ext == ".zip" {
							cmd = exec.Command("unzip", target.Path, "-d", m.CurrentDir)
						} else if ext == ".tar" || ext == ".tar.gz" || ext == ".tgz" {
							cmd = exec.Command("tar", "-xvf", target.Path, "-C", m.CurrentDir)
						}
						
						if cmd != nil {
							err := cmd.Run()
							if err != nil {
								m.LastCmdOut = m.Styles.ErrorText.Render(fmt.Sprintf("✗ Extraction failed: %v", err))
							} else {
								m.LastCmdOut = m.Styles.SuccessText.Render(fmt.Sprintf("✓ Extracted %s successfully", target.Name))
								m.refreshLeftEntries()
								m.refreshCenterAndRight()
							}
						} else {
							m.StatusMsg = "Not a supported archive (.zip, .tar, .tar.gz)"
						}
					}
				}
				return m, nil

			case "o", "O":
				if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if !target.IsDir {
						var cmd *exec.Cmd
						macCheck := exec.Command("uname")
						out, _ := macCheck.Output()
						if strings.Contains(strings.ToLower(string(out)), "darwin") {
							cmd = exec.Command("open", target.Path)
						} else {
							cmd = exec.Command("xdg-open", target.Path)
						}
						if cmd != nil {
							_ = cmd.Start()
							m.StatusMsg = fmt.Sprintf("Opened: %s", target.Name)
						}
					}
				}
				return m, nil

			case "b":
				m.BookmarkCursor = 0
				m.BookmarkScroll = 0
				m.State = StateBookmarks
				return m, nil

			case "B":
				m.Bookmarks = config.AddBookmark(m.Bookmarks, m.CurrentDir)
				config.SaveBookmarks(m.Bookmarks)
				m.StatusMsg = fmt.Sprintf("★ Bookmarked: %s", m.CurrentDir)
				return m, nil

			case "g":
				switch m.Focus {
				case PaneLeft:
					m.LeftCursor = 0
					m.LeftScroll = 0
					m.RightScroll = 0
					m.refreshCenterAndRight()
				case PaneCenter:
					m.CenterCursor = 0
					m.CenterScroll = 0
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
				return m, nil

			case "G":
				switch m.Focus {
				case PaneLeft:
					m.LeftCursor = max(0, len(m.LeftEntries)-1)
					m.RightScroll = 0
					m.refreshCenterAndRight()
				case PaneCenter:
					m.CenterCursor = max(0, len(m.CenterEntries)-1)
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
				return m, nil
			}
		} else {
			// Terminal pane focused
			switch msg.String() {
			case "esc":
				m.Focus = PaneLeft
				m.TermInput.Blur()
				m.HistoryIdx = -1
				return m, nil

			case "up":
				if len(m.CmdHistory) > 0 {
					if m.HistoryIdx == -1 {
						m.HistorySaved = m.TermInput.Value()
					}
					if m.HistoryIdx < len(m.CmdHistory)-1 {
						m.HistoryIdx++
						m.TermInput.SetValue(m.CmdHistory[m.HistoryIdx])
						m.TermInput.CursorEnd()
					}
				}
				return m, nil

			case "down":
				if m.HistoryIdx > 0 {
					m.HistoryIdx--
					m.TermInput.SetValue(m.CmdHistory[m.HistoryIdx])
					m.TermInput.CursorEnd()
				} else if m.HistoryIdx == 0 {
					m.HistoryIdx = -1
					m.TermInput.SetValue(m.HistorySaved)
					m.TermInput.CursorEnd()
				}
				return m, nil

			case "enter":
				cmdStr := m.TermInput.Value()
				m.TermInput.SetValue("")
				m.HistoryIdx = -1
				m.HistorySaved = ""

				if strings.TrimSpace(cmdStr) != "" {
					m.CmdHistory = config.AddHistory(m.CmdHistory, strings.TrimSpace(cmdStr))
					config.SaveHistory(m.CmdHistory)

					trimmed := strings.TrimSpace(cmdStr)
					if strings.HasPrefix(trimmed, "rm ") || strings.HasPrefix(trimmed, "chmod ") || strings.HasPrefix(trimmed, "chown ") {
						m.Dialog = ui.NewPermissionDialog(
							"Permission Warning",
							fmt.Sprintf("Allow execution of destructive command: '%s'?", trimmed),
							func() {
								res := terminal.ExecuteCommand(trimmed, m.CurrentDir, m.Config)
								m.handleCommandResult(res)
							},
						)
						m.State = StateDialog
						return m, nil
					}

					res := terminal.ExecuteCommand(cmdStr, m.CurrentDir, m.Config)
					m.handleCommandResult(res)
				}
				return m, nil
			}
		}

		// Navigation inside active focus pane
		switch m.Focus {
		case PaneLeft:
			switch msg.String() {
			case "up", "k":
				if m.LeftCursor > 0 {
					m.LeftCursor--
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
			case "down", "j":
				if m.LeftCursor < len(m.LeftEntries)-1 {
					m.LeftCursor++
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
			case "right", "l", "enter":
				if len(m.LeftEntries) > 0 {
					selected := m.LeftEntries[m.LeftCursor]
					if selected.IsDir {
						if selected.Name == ".." {
							parentDir := filepath.Dir(m.CurrentDir)
							if parentDir != m.CurrentDir {
								m.CurrentDir = parentDir
								m.LeftCursor = 0
								m.LeftScroll = 0
								m.CenterCursor = 0
								m.CenterScroll = 0
								m.RightScroll = 0
								m.FilterInput.SetValue("")
								m.refreshLeftEntries()
								m.refreshCenterAndRight()
							}
						} else {
							m.Focus = PaneCenter
							m.CenterCursor = 0
							m.CenterScroll = 0
							m.RightScroll = 0
							m.refreshCenterAndRight()
						}
					}
				}
			case "left", "h":
				parentDir := filepath.Dir(m.CurrentDir)
				if parentDir != m.CurrentDir {
					m.CurrentDir = parentDir
					m.LeftCursor = 0
					m.LeftScroll = 0
					m.CenterCursor = 0
					m.CenterScroll = 0
					m.RightScroll = 0
					m.FilterInput.SetValue("")
					m.refreshLeftEntries()
					m.refreshCenterAndRight()
				}
			}

		case PaneCenter:
			switch msg.String() {
			case "up", "k":
				if m.CenterCursor > 0 {
					m.CenterCursor--
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
			case "down", "j":
				if m.CenterCursor < len(m.CenterEntries)-1 {
					m.CenterCursor++
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
			case "right", "l", "enter":
				if len(m.CenterEntries) > 0 {
					selected := m.CenterEntries[m.CenterCursor]
					if selected.IsDir {
						m.CurrentDir = selected.Path
						m.LeftCursor = 0
						m.LeftScroll = 0
						m.CenterCursor = 0
						m.CenterScroll = 0
						m.RightScroll = 0
						m.FilterInput.SetValue("")
						m.refreshLeftEntries()
						m.refreshCenterAndRight()
					}
				}
			case "left", "h":
				m.RightScroll = 0
				m.refreshCenterAndRight()
				m.Focus = PaneLeft
			}

		case PaneRight:
			switch msg.String() {
			case "up", "k":
				if m.RightScroll > 0 {
					m.RightScroll--
				}
			case "down", "j":
				m.RightScroll++
			case "left", "h":
				m.Focus = PaneLeft
			}
		}
	}

	if m.Focus == PaneTerminal {
		m.TermInput, cmd = m.TermInput.Update(msg)
	}

	return m, cmd
}

func (m *AppModel) handleCommandResult(res terminal.CommandResult) {
	if res.Err != nil && res.Output != "" {
		m.LastCmdOut = m.Styles.ErrorText.Render("✗ ") + res.Output
	} else {
		m.LastCmdOut = res.Output
	}
	if res.NewDir != "" {
		m.CurrentDir = res.NewDir
		m.LeftCursor = 0
		m.LeftScroll = 0
		m.CenterCursor = 0
		m.CenterScroll = 0
		m.FilterInput.SetValue("")
		m.refreshLeftEntries()
		m.refreshCenterAndRight()
		m.StatusMsg = fmt.Sprintf("Directory changed to: %s", m.CurrentDir)
	} else {
		m.refreshLeftEntries()
		m.refreshCenterAndRight()
	}
}

// renderFileItem renders a single file list row with zebra-striping, size column, accent bar.
func (m AppModel) renderFileItem(entry filetree.FileEntry, idx int, cursor int, isFocused bool, paneWidth int) string {
	icon := entry.Icon
	name := entry.Name
	if entry.IsDir && name != ".." {
		name = name + "/"
	}

	sizeStr := ""
	if !entry.IsDir {
		sizeStr = entry.FormatSize()
	}

	// Multi-select prefix
	selectPrefix := "  "
	if m.MultiSelect[entry.Path] {
		selectPrefix = "✓ "
	}

	// Git status styling
	gitSuffix := ""
	unrenderedGitSuffix := ""
	if entry.GitStatus != "" {
		unrenderedGitSuffix = " [" + entry.GitStatus + "]"
		switch entry.GitStatus {
		case "M":
			gitSuffix = m.Styles.GitModified.Render(" [M]")
		case "A":
			gitSuffix = m.Styles.GitAdded.Render(" [A]")
		case "?":
			gitSuffix = m.Styles.GitUntracked.Render(" [?]")
		default:
			gitSuffix = m.Styles.MutedText.Render(fmt.Sprintf(" [%s]", entry.GitStatus))
		}
	}

	// Build the item text
	nameDisplay := fmt.Sprintf("%s%s %s", selectPrefix, icon, name)
	if entry.IsSymlink && entry.SymlinkTarget != "" {
		nameDisplay += fmt.Sprintf(" → %s", entry.SymlinkTarget)
	}

	// Calculate available width for name (leave room for size)
	availWidth := max(10, paneWidth-len(sizeStr)-6)
	unrenderedTotal := len(nameDisplay) + len(unrenderedGitSuffix)
	if unrenderedTotal > availWidth {
		allowedNameLen := availWidth - len(unrenderedGitSuffix)
		if len(nameDisplay) > allowedNameLen && allowedNameLen > 3 {
			nameDisplay = nameDisplay[:allowedNameLen-1] + "…"
		}
	}

	// Pad to fill width
	gap := max(0, paneWidth-len(nameDisplay)-len(unrenderedGitSuffix)-len(sizeStr)-4)
	line := nameDisplay + gitSuffix + strings.Repeat(" ", gap) + sizeStr

	if idx == cursor {
		prefix := "  "
		if isFocused {
			prefix = m.Styles.AccentBar.Render("▌") + " "
		}
		return m.Styles.SelectedItem.Render(prefix + line)
	}

	prefix := "  "
	isZebra := idx%2 == 1

	if entry.IsDir {
		if isZebra {
			return m.Styles.ZebraDirItem.Render(prefix + line)
		}
		return m.Styles.DirectoryItem.Render(prefix + line)
	}
	if isZebra {
		return m.Styles.ZebraItem.Render(prefix + line)
	}
	return m.Styles.Item.Render(prefix + line)
}

func (m AppModel) View() string {
	if m.State == StateStartup {
		return m.Startup.View()
	}

	if m.State == StateAnalyzer || m.State == StateAnalyzerLoading {
		return m.renderAnalyzerView()
	}

	if m.State == StateBookmarks {
		return m.renderBookmarksView()
	}

	if m.Width <= 0 || m.Height <= 0 {
		return "Initializing Matt..."
	}

	var viewStr string

	// 1. Header Bar with breadcrumb
	headerLeft := m.Styles.HeaderTitle.Render(" MATT ") + " " + m.buildBreadcrumb()
	headerRight := m.getModeIndicator() + " " + m.Styles.HeaderBadge.Render(fmt.Sprintf(" %s ", version.Version))
	headerGap := max(0, m.Width-lipgloss.Width(headerLeft)-lipgloss.Width(headerRight))
	headerBar := m.Styles.Header.Render(headerLeft + strings.Repeat(" ", headerGap) + headerRight)

	// 2. Main 3-Column Height & Width Calculations
	mainHeight := max(10, m.Height-10)
	leftWidth := max(22, m.Width/4)
	centerWidth := max(22, m.Width/4)
	rightWidth := max(30, m.Width-leftWidth-centerWidth-6)

	maxVisibleRows := max(1, mainHeight-4) // -4 for title + scroll indicators

	// --- Left Pane ---
	leftTitle := "Navigation"
	if m.IsFiltering || m.FilterInput.Value() != "" {
		leftTitle = fmt.Sprintf("Nav [%s]", m.FilterInput.View())
	}
	leftTitleStyle := m.Styles.PaneTitleInactive
	leftBorderStyle := m.Styles.InactiveBorder
	focusIcon := "○"
	if m.Focus == PaneLeft {
		leftTitleStyle = m.Styles.PaneTitleActive
		leftBorderStyle = m.Styles.ActiveBorder
		focusIcon = "◉"
	}

	// Update Left Scroll Window
	if m.LeftCursor < m.LeftScroll {
		m.LeftScroll = m.LeftCursor
	}
	if m.LeftCursor >= m.LeftScroll+maxVisibleRows {
		m.LeftScroll = m.LeftCursor - maxVisibleRows + 1
	}
	if m.LeftScroll < 0 {
		m.LeftScroll = 0
	}

	lDirs, lFiles := countDirsFiles(m.LeftEntries)
	var leftItems []string
	
	// Format and pad left header to exactly fit leftWidth (excluding borders)
	leftHeaderStr := fmt.Sprintf(" %s %s [%d/%d] ", focusIcon, leftTitle, m.LeftCursor+1, len(m.LeftEntries))
	leftCountsStr := fmt.Sprintf("%dd %df ", lDirs, lFiles)
	if len(leftHeaderStr)+len(leftCountsStr) <= leftWidth {
		leftHeaderStr += leftCountsStr
	}
	if len(leftHeaderStr) > leftWidth {
		leftHeaderStr = leftHeaderStr[:leftWidth]
	} else {
		leftHeaderStr += strings.Repeat(" ", max(0, leftWidth-len(leftHeaderStr)))
	}
	leftItems = append(leftItems, leftTitleStyle.Render(leftHeaderStr))

	// Top scroll indicator
	if m.LeftScroll > 0 {
		leftItems = append(leftItems, m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▲ %d more above", m.LeftScroll)))
	}

	lStartIdx := m.LeftScroll
	if lStartIdx >= len(m.LeftEntries) {
		lStartIdx = 0
	}
	lEndIdx := min(len(m.LeftEntries), lStartIdx+maxVisibleRows)

	for i := lStartIdx; i < lEndIdx; i++ {
		leftItems = append(leftItems, m.renderFileItem(m.LeftEntries[i], i, m.LeftCursor, m.Focus == PaneLeft, leftWidth))
	}

	// Bottom scroll indicator
	if lEndIdx < len(m.LeftEntries) {
		leftItems = append(leftItems, m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▼ %d more below", len(m.LeftEntries)-lEndIdx)))
	}

	leftContent := strings.Join(leftItems, "\n")
	leftBox := leftBorderStyle.Width(leftWidth).Height(mainHeight).Render(leftContent)

	// --- Center Pane ---
	centerTitle := "Explorer"
	centerTitleStyle := m.Styles.PaneTitleInactive
	centerBorderStyle := m.Styles.InactiveBorder
	centerFocusIcon := "○"
	if m.Focus == PaneCenter {
		centerTitleStyle = m.Styles.PaneTitleActive
		centerBorderStyle = m.Styles.ActiveBorder
		centerFocusIcon = "◉"
	}

	// Update Center Scroll Window
	if m.CenterCursor < m.CenterScroll {
		m.CenterScroll = m.CenterCursor
	}
	if m.CenterCursor >= m.CenterScroll+maxVisibleRows {
		m.CenterScroll = m.CenterCursor - maxVisibleRows + 1
	}
	if m.CenterScroll < 0 {
		m.CenterScroll = 0
	}

	cDirs, cFiles := countDirsFiles(m.CenterEntries)
	var centerItems []string
	
	// Format and pad center header to exactly fit centerWidth
	var centerHeaderStr string
	if len(m.CenterEntries) > 0 {
		centerHeaderStr = fmt.Sprintf(" %s %s [%d/%d] ", centerFocusIcon, centerTitle, m.CenterCursor+1, len(m.CenterEntries))
		centerCountsStr := fmt.Sprintf("%dd %df ", cDirs, cFiles)
		if len(centerHeaderStr)+len(centerCountsStr) <= centerWidth {
			centerHeaderStr += centerCountsStr
		}
	} else {
		centerHeaderStr = fmt.Sprintf(" %s %s (empty) ", centerFocusIcon, centerTitle)
	}
	if len(centerHeaderStr) > centerWidth {
		centerHeaderStr = centerHeaderStr[:centerWidth]
	} else {
		centerHeaderStr += strings.Repeat(" ", max(0, centerWidth-len(centerHeaderStr)))
	}
	centerItems = append(centerItems, centerTitleStyle.Render(centerHeaderStr))

	// Top scroll indicator
	if m.CenterScroll > 0 {
		centerItems = append(centerItems, m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▲ %d more above", m.CenterScroll)))
	}

	cStartIdx := m.CenterScroll
	if cStartIdx >= len(m.CenterEntries) {
		cStartIdx = 0
	}
	cEndIdx := min(len(m.CenterEntries), cStartIdx+maxVisibleRows)

	for i := cStartIdx; i < cEndIdx; i++ {
		centerItems = append(centerItems, m.renderFileItem(m.CenterEntries[i], i, m.CenterCursor, m.Focus == PaneCenter, centerWidth))
	}

	// Bottom scroll indicator
	if cEndIdx < len(m.CenterEntries) {
		centerItems = append(centerItems, m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▼ %d more below", len(m.CenterEntries)-cEndIdx)))
	}

	centerContent := strings.Join(centerItems, "\n")
	centerBox := centerBorderStyle.Width(centerWidth).Height(mainHeight).Render(centerContent)

	// --- Right Column: Top Preview Box & Bottom Metadata Box ---
	metaHeight := 7
	previewHeight := max(4, mainHeight-metaHeight)

	rightTitleStyle := m.Styles.PaneTitleInactive
	rightBorderStyle := m.Styles.InactiveBorder
	rightFocusIcon := "○"
	if m.Focus == PaneRight {
		rightTitleStyle = m.Styles.PaneTitleActive
		rightBorderStyle = m.Styles.ActiveBorder
		rightFocusIcon = "◉"
	}

	// 1. Top-Right Preview Box
	var rightItems []string
	
	// Format and pad right header to exactly fit rightWidth
	rightHeaderStr := fmt.Sprintf(" %s Preview [%s] ", rightFocusIcon, m.RightPreview.Info)
	if len(rightHeaderStr) > rightWidth {
		rightHeaderStr = rightHeaderStr[:rightWidth]
	} else {
		rightHeaderStr += strings.Repeat(" ", max(0, rightWidth-len(rightHeaderStr)))
	}
	rightItems = append(rightItems, rightTitleStyle.Render(rightHeaderStr))

	previewLines := strings.Split(m.RightPreview.Content, "\n")
	rStartIdx := max(0, m.RightScroll)
	rEndIdx := min(len(previewLines), rStartIdx+(previewHeight-3))
	if rStartIdx > 0 {
		rightItems = append(rightItems, m.Styles.ScrollIndicator.Render("  ▲ scroll up"))
	}
	for i := rStartIdx; i < rEndIdx; i++ {
		rightItems = append(rightItems, previewLines[i])
	}
	if rEndIdx < len(previewLines) {
		rightItems = append(rightItems, m.Styles.ScrollIndicator.Render("  ▼ scroll down"))
	}
	rightContent := strings.Join(rightItems, "\n")
	rightPreviewBox := rightBorderStyle.Width(rightWidth).Height(previewHeight).Render(rightContent)

	// 2. Bottom-Right Metadata Inspector Box
	var metaSb strings.Builder
	metaSb.WriteString(m.Styles.PaneTitleInactive.Render(" File Metadata Inspector "))
	metaSb.WriteString("\n")

	sel := m.SelectedItem
	if sel.Name != "" {
		metaSb.WriteString(fmt.Sprintf("  %s %s\n", m.Styles.MetaLabel.Render("Name:"), m.Styles.MetaValue.Render(sel.Name)))
		metaSb.WriteString(fmt.Sprintf("  %s %s  %s %s\n", m.Styles.MetaLabel.Render("Type:"), m.Styles.MetaValue.Render(sel.Extension), m.Styles.MetaLabel.Render("Size:"), m.Styles.MetaValue.Render(sel.FormatSize())))
		metaSb.WriteString(fmt.Sprintf("  %s %s  %s %s:%s\n", m.Styles.MetaLabel.Render("Perms:"), m.Styles.MetaValue.Render(sel.Permissions), m.Styles.MetaLabel.Render("Owner:"), m.Styles.MetaValue.Render(sel.Owner), m.Styles.MetaValue.Render(sel.Group)))
		metaSb.WriteString(fmt.Sprintf("  %s %s\n", m.Styles.MetaLabel.Render("Modified:"), m.Styles.MetaValue.Render(sel.FormatModTime())))
		if sel.IsSymlink && sel.SymlinkTarget != "" {
			metaSb.WriteString(fmt.Sprintf("  %s %s\n", m.Styles.MetaLabel.Render("Link:"), m.Styles.MetaValue.Render(sel.SymlinkTarget)))
		}
	} else {
		metaSb.WriteString("  [No item selected]")
	}

	rightMetaBox := m.Styles.InactiveBorder.Width(rightWidth).Height(metaHeight).Render(metaSb.String())

	// Combine Top-Right Preview Box and Bottom-Right Metadata Box vertically
	rightColumn := lipgloss.JoinVertical(lipgloss.Left, rightPreviewBox, rightMetaBox)

	// Combine 3 Upper Panes horizontally
	upperPanes := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, centerBox, rightColumn)

	// --- Bottom Terminal Pane ---
	termTitleStyle := m.Styles.PaneTitleInactive
	termBorderStyle := m.Styles.InactiveBorder
	if m.Focus == PaneTerminal {
		termTitleStyle = m.Styles.PaneTitleActive
		termBorderStyle = m.Styles.ActiveBorder
	}

	var termSb strings.Builder
	termSb.WriteString(termTitleStyle.Render(" Terminal [':' to focus] "))
	if len(m.ClipboardPaths) > 0 {
		clipboardDesc := filepath.Base(m.ClipboardPaths[0])
		if len(m.ClipboardPaths) > 1 {
			clipboardDesc = fmt.Sprintf("%d items", len(m.ClipboardPaths))
		}
		termSb.WriteString("  " + m.Styles.WarningText.Render(fmt.Sprintf("📋 %s: %s", m.ClipboardOp, clipboardDesc)))
	}
	termSb.WriteString("\n")
	if m.LastCmdOut != "" {
		outLines := strings.Split(m.LastCmdOut, "\n")
		if len(outLines) > 2 {
			outLines = outLines[len(outLines)-2:]
		}
		termSb.WriteString(m.Styles.MutedText.Render(fmt.Sprintf("  %s\n", strings.Join(outLines, " | "))))
	}

	// Input prompt overlay
	if m.State == StateInput {
		var label string
		switch m.InputAction {
		case InputNewFile:
			label = "New File: "
		case InputNewFolder:
			label = "New Folder: "
		case InputRename:
			label = "Rename: "
		}
		termSb.WriteString("  " + m.Styles.WarningText.Render(label) + m.PromptInput.View())
	} else {
		termSb.WriteString("  " + m.TermInput.View())
	}

	termBox := termBorderStyle.Width(m.Width - 2).Height(5).Render(termSb.String())

	// --- Bottom Status Bar ---
	statusLeft := " [Tab] Pane • [/] Filter • [Alt+D] Analyzer • [n/N] New • [c/x/p] Copy/Cut/Paste • [b/B] Bookmarks • [g/G] Jump"
	statusRight := m.StatusMsg
	statusGap := max(0, m.Width-lipgloss.Width(statusLeft)-lipgloss.Width(statusRight))
	statusBar := m.Styles.StatusBar.Render(statusLeft + strings.Repeat(" ", statusGap) + statusRight)

	viewStr = lipgloss.JoinVertical(
		lipgloss.Left,
		headerBar,
		upperPanes,
		termBox,
		statusBar,
	)

	rootStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.Config.Theme.Bg)).
		Width(m.Width).
		Height(m.Height)

	if m.State == StateDialog {
		return rootStyle.Render(m.Dialog.View(m.Styles, m.Width, m.Height))
	}

	return rootStyle.Render(viewStr)
}

func (m AppModel) renderAnalyzerView() string {
	rootStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.Config.Theme.Bg)).
		Width(m.Width).
		Height(m.Height)

	if m.State == StateAnalyzerLoading {
		loadingContent := m.Styles.ModalTitle.Render("Disk Space Analyzer") + "\n\n" +
			m.Styles.MutedText.Render("  Scanning directory...\n  Please wait while disk usage is analyzed.\n\n") +
			m.Styles.MutedText.Render("  Path: "+m.CurrentDir)

		box := m.Styles.ModalBox.Render(loadingContent)
		placed := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box, lipgloss.WithWhitespaceChars(" "))
		return rootStyle.Render(placed)
	}

	var sb strings.Builder
	sb.WriteString(m.Styles.ModalTitle.Render(fmt.Sprintf("Disk Space Analyzer: %s", m.AnalyzerResult.DirPath)))
	sb.WriteString("\n")

	sortLabel := "Size ↓"
	if m.AnalyzerResult.SortByName {
		sortLabel = "Name A-Z"
	}
	sb.WriteString(m.Styles.MutedText.Render(fmt.Sprintf("   Total: %s  •  Items: %d  •  Sort: %s",
		analyzer.FormatBytes(m.AnalyzerResult.TotalSize), len(m.AnalyzerResult.Items), sortLabel)))
	sb.WriteString("\n\n")

	maxRows := max(5, m.Height-14)
	items := m.AnalyzerResult.Items

	// Scroll logic for analyzer
	if m.AnalyzerCursor < m.AnalyzerScroll {
		m.AnalyzerScroll = m.AnalyzerCursor
	}
	if m.AnalyzerCursor >= m.AnalyzerScroll+maxRows {
		m.AnalyzerScroll = m.AnalyzerCursor - maxRows + 1
	}

	startIdx := m.AnalyzerScroll
	if startIdx >= len(items) {
		startIdx = 0
	}
	endIdx := min(len(items), startIdx+maxRows)

	if startIdx > 0 {
		sb.WriteString(m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▲ %d items above", startIdx)))
		sb.WriteString("\n")
	}

	for i := startIdx; i < endIdx; i++ {
		item := items[i]
		prefix := "DIR "
		if !item.IsDir {
			prefix = "FILE"
		}
		line := fmt.Sprintf("  %s %-20s  %s  %5.1f%%  %-10s", prefix, item.Name, item.Bar, item.Percentage, analyzer.FormatBytes(item.SizeBytes))
		if i == m.AnalyzerCursor {
			sb.WriteString(m.Styles.SelectedItem.Render(m.Styles.AccentBar.Render("▌") + " " + line))
		} else {
			if i%2 == 1 {
				sb.WriteString(m.Styles.ZebraItem.Render("  " + line))
			} else {
				sb.WriteString(m.Styles.Item.Render("  " + line))
			}
		}
		sb.WriteString("\n")
	}

	if endIdx < len(items) {
		sb.WriteString(m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▼ %d items below", len(items)-endIdx)))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(m.Styles.MutedText.Render(" [↑/↓] Navigate  •  [Enter] Open Dir  •  [s] Toggle Sort  •  [Alt+D / Esc] Close"))

	box := m.Styles.ModalBox.Render(sb.String())
	placed := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars(" "),
	)
	return rootStyle.Render(placed)
}

func (m AppModel) renderBookmarksView() string {
	rootStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.Config.Theme.Bg)).
		Width(m.Width).
		Height(m.Height)

	var sb strings.Builder
	sb.WriteString(m.Styles.ModalTitle.Render("Bookmarks"))
	sb.WriteString("\n\n")

	if len(m.Bookmarks) == 0 {
		sb.WriteString(m.Styles.MutedText.Render("  No bookmarks saved yet.\n  Press [B] in main view to bookmark current directory."))
	} else {
		maxRows := max(5, m.Height-12)
		startIdx := m.BookmarkScroll
		if m.BookmarkCursor < startIdx {
			startIdx = m.BookmarkCursor
		}
		if m.BookmarkCursor >= startIdx+maxRows {
			startIdx = m.BookmarkCursor - maxRows + 1
		}
		m.BookmarkScroll = startIdx

		endIdx := min(len(m.Bookmarks), startIdx+maxRows)

		if startIdx > 0 {
			sb.WriteString(m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▲ %d more above", startIdx)))
			sb.WriteString("\n")
		}

		for i := startIdx; i < endIdx; i++ {
			bm := m.Bookmarks[i]
			// Shorten home directory to ~
			display := bm
			home := config.GetHomeDir()
			if strings.HasPrefix(display, home) {
				display = "~" + display[len(home):]
			}

			if i == m.BookmarkCursor {
				sb.WriteString(m.Styles.SelectedItem.Render(fmt.Sprintf("%s 📁 %s", m.Styles.AccentBar.Render("▌"), display)))
			} else {
				if i%2 == 1 {
					sb.WriteString(m.Styles.ZebraItem.Render(fmt.Sprintf("  📁 %s", display)))
				} else {
					sb.WriteString(m.Styles.Item.Render(fmt.Sprintf("  📁 %s", display)))
				}
			}
			sb.WriteString("\n")
		}

		if endIdx < len(m.Bookmarks) {
			sb.WriteString(m.Styles.ScrollIndicator.Render(fmt.Sprintf("  ▼ %d more below", len(m.Bookmarks)-endIdx)))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(m.Styles.MutedText.Render(" [↑/↓] Navigate  •  [Enter] Jump  •  [d] Remove  •  [Esc] Close"))

	box := m.Styles.ModalBox.Render(sb.String())
	placed := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars(" "),
	)
	return rootStyle.Render(placed)
}
