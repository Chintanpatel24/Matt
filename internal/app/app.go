package app

import (
	"fmt"
	"os"
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
)

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
	TermInput   textinput.Model
	FilterInput textinput.Model
	IsFiltering  bool

	// Analyzer state
	AnalyzerResult analyzer.DiskUsageResult
	AnalyzerCursor int

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
	ti.Placeholder = "Type shell command or extension alias (e.g. cd, touch, findbig, ll)..."
	ti.Prompt = "matt $ "
	ti.CharLimit = 256
	ti.Width = 60

	fi := textinput.New()
	fi.Placeholder = "Fuzzy search files..."
	fi.Prompt = "🔍 / "
	fi.CharLimit = 64
	fi.Width = 24

	styles := ui.NewStyles(cfg)
	startup := NewStartupModel(cfg, absDir)

	m := AppModel{
		State:        StateStartup,
		Config:       cfg,
		Styles:       styles,
		CurrentDir:   absDir,
		Focus:        PaneLeft,
		TermInput:    ti,
		FilterInput:  fi,
		IsFiltering:  false,
		Startup:      startup,
		LastCmdOut:   "Ready. Press 'Tab' to switch focus, ':' for terminal, '/' for search, 'Alt+D' for Analyzer.",
		StatusMsg:    "Matt File Manager - High Performance Matt Black TUI",
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
			m.CenterCursor = 0
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

	case tea.MouseMsg:
		if m.State == StateMain && msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			leftWidth := m.Width / 4
			if msg.X >= 1 && msg.X <= leftWidth {
				m.Focus = PaneLeft
				clickedRow := msg.Y - 3
				if clickedRow >= 0 && clickedRow < len(m.LeftEntries) {
					m.LeftCursor = clickedRow
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
			} else if msg.X > leftWidth && msg.X <= leftWidth*2 {
				m.Focus = PaneCenter
				clickedRow := msg.Y - 3
				if clickedRow >= 0 && clickedRow < len(m.CenterEntries) {
					m.CenterCursor = clickedRow
					m.RightScroll = 0
					m.refreshCenterAndRight()
				}
			} else if msg.X > leftWidth*2 && msg.X < m.Width {
				m.Focus = PaneRight
			}
		}
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

	// 3. Disk Space Analyzer View Mode
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

	// 4. Main View Mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "alt+d" || msg.String() == "alt+D" {
			res, err := analyzer.AnalyzeDirectory(m.CurrentDir)
			if err == nil {
				m.AnalyzerResult = res
				m.AnalyzerCursor = 0
				m.State = StateAnalyzer
				return m, nil
			}
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
				return m, nil

			case "r":
				m.refreshLeftEntries()
				m.refreshCenterAndRight()
				m.StatusMsg = "Refreshed view."
				return m, nil

			case "d":
				if len(m.LeftEntries) > 0 {
					target := m.LeftEntries[m.LeftCursor]
					if target.Name != ".." {
						m.Dialog = ui.NewPermissionDialog(
							"Delete Confirmation",
							fmt.Sprintf("Are you sure you want to delete '%s'?", target.Name),
							func() {
								err := os.RemoveAll(target.Path)
								if err != nil {
									m.LastCmdOut = fmt.Sprintf("Error deleting %s: %v", target.Name, err)
								} else {
									m.LastCmdOut = fmt.Sprintf("Deleted '%s' successfully.", target.Name)
									m.refreshLeftEntries()
									m.refreshCenterAndRight()
								}
							},
						)
						m.State = StateDialog
						return m, nil
					}
				}
			}
		} else {
			switch msg.String() {
			case "esc":
				m.Focus = PaneLeft
				m.TermInput.Blur()
				return m, nil

			case "enter":
				cmdStr := m.TermInput.Value()
				m.TermInput.SetValue("")

				if strings.TrimSpace(cmdStr) != "" {
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
						m.Focus = PaneCenter
						m.CenterCursor = 0
						m.CenterScroll = 0
						m.RightScroll = 0
						m.refreshCenterAndRight()
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
	m.LastCmdOut = res.Output
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

func (m AppModel) View() string {
	if m.State == StateStartup {
		return m.Startup.View()
	}

	if m.State == StateAnalyzer {
		return m.renderAnalyzerView()
	}

	if m.Width <= 0 || m.Height <= 0 {
		return "Initializing Matt..."
	}

	var viewStr string

	// 1. Header Bar
	headerLeft := m.Styles.HeaderTitle.Render(" MATT BLACK ") + " " + m.Styles.HeaderPath.Render(m.CurrentDir)
	headerRight := m.Styles.HeaderBadge.Render(" Matt TUI v1.0.0 ")
	headerGap := max(0, m.Width-lipgloss.Width(headerLeft)-lipgloss.Width(headerRight))
	headerBar := m.Styles.Header.Render(headerLeft + strings.Repeat(" ", headerGap) + headerRight)

	// 2. Main 3-Column Height & Width Calculations
	mainHeight := max(10, m.Height-10)
	leftWidth := max(22, m.Width/4)
	centerWidth := max(22, m.Width/4)
	rightWidth := max(30, m.Width-leftWidth-centerWidth-6)

	maxVisibleRows := max(1, mainHeight-3)

	// --- Left Pane (Navigation Robust Windowed Viewport) ---
	leftTitle := "Navigation"
	if m.IsFiltering || m.FilterInput.Value() != "" {
		leftTitle = fmt.Sprintf("Nav [%s]", m.FilterInput.View())
	}
	leftTitleStyle := m.Styles.PaneTitleInactive
	leftBorderStyle := m.Styles.InactiveBorder
	if m.Focus == PaneLeft {
		leftTitleStyle = m.Styles.PaneTitleActive
		leftBorderStyle = m.Styles.ActiveBorder
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

	var leftItems []string
	leftItems = append(leftItems, leftTitleStyle.Render(fmt.Sprintf(" %s (%d) ", leftTitle, len(m.LeftEntries))))

	lStartIdx := m.LeftScroll
	if lStartIdx >= len(m.LeftEntries) {
		lStartIdx = 0
	}
	lEndIdx := min(len(m.LeftEntries), lStartIdx+maxVisibleRows)

	for i := lStartIdx; i < lEndIdx; i++ {
		entry := m.LeftEntries[i]
		itemText := fmt.Sprintf("%s %s", entry.Icon, entry.Name)
		if entry.IsDir {
			itemText = fmt.Sprintf("%s %s/", entry.Icon, entry.Name)
		}

		if i == m.LeftCursor {
			if m.Focus == PaneLeft {
				itemText = m.Styles.SelectedItem.Render(fmt.Sprintf("> %s", itemText))
			} else {
				itemText = m.Styles.SelectedItem.Render(fmt.Sprintf("  %s", itemText))
			}
		} else {
			if entry.IsDir {
				itemText = m.Styles.DirectoryItem.Render(fmt.Sprintf("  %s", itemText))
			} else {
				itemText = m.Styles.Item.Render(fmt.Sprintf("  %s", itemText))
			}
		}
		leftItems = append(leftItems, itemText)
	}
	leftContent := strings.Join(leftItems, "\n")
	leftBox := leftBorderStyle.Width(leftWidth).Height(mainHeight).Render(leftContent)

	// --- Center Pane (Explorer Robust Windowed Viewport) ---
	centerTitle := "Explorer"
	centerTitleStyle := m.Styles.PaneTitleInactive
	centerBorderStyle := m.Styles.InactiveBorder
	if m.Focus == PaneCenter {
		centerTitleStyle = m.Styles.PaneTitleActive
		centerBorderStyle = m.Styles.ActiveBorder
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

	var centerItems []string
	centerItems = append(centerItems, centerTitleStyle.Render(fmt.Sprintf(" %s (%d) ", centerTitle, len(m.CenterEntries))))

	cStartIdx := m.CenterScroll
	if cStartIdx >= len(m.CenterEntries) {
		cStartIdx = 0
	}
	cEndIdx := min(len(m.CenterEntries), cStartIdx+maxVisibleRows)

	for i := cStartIdx; i < cEndIdx; i++ {
		entry := m.CenterEntries[i]
		itemText := fmt.Sprintf("%s %s", entry.Icon, entry.Name)
		if entry.IsDir {
			itemText = fmt.Sprintf("%s %s/", entry.Icon, entry.Name)
		}

		if i == m.CenterCursor {
			if m.Focus == PaneCenter {
				itemText = m.Styles.SelectedItem.Render(fmt.Sprintf("> %s", itemText))
			} else {
				itemText = m.Styles.SelectedItem.Render(fmt.Sprintf("  %s", itemText))
			}
		} else {
			if entry.IsDir {
				itemText = m.Styles.DirectoryItem.Render(fmt.Sprintf("  %s", itemText))
			} else {
				itemText = m.Styles.Item.Render(fmt.Sprintf("  %s", itemText))
			}
		}
		centerItems = append(centerItems, itemText)
	}
	centerContent := strings.Join(centerItems, "\n")
	centerBox := centerBorderStyle.Width(centerWidth).Height(mainHeight).Render(centerContent)

	// --- Right Column: Top Preview Box & Bottom Metadata Box ---
	metaHeight := 7
	previewHeight := max(4, mainHeight-metaHeight)

	rightTitleStyle := m.Styles.PaneTitleInactive
	rightBorderStyle := m.Styles.InactiveBorder
	if m.Focus == PaneRight {
		rightTitleStyle = m.Styles.PaneTitleActive
		rightBorderStyle = m.Styles.ActiveBorder
	}

	// 1. Top-Right Preview Box
	var rightItems []string
	rightItems = append(rightItems, rightTitleStyle.Render(fmt.Sprintf(" Preview [%s] ", m.RightPreview.Info)))

	previewLines := strings.Split(m.RightPreview.Content, "\n")
	rStartIdx := max(0, m.RightScroll)
	rEndIdx := min(len(previewLines), rStartIdx+(previewHeight-3))
	for i := rStartIdx; i < rEndIdx; i++ {
		rightItems = append(rightItems, previewLines[i])
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
	termSb.WriteString(termTitleStyle.Render(" Terminal & Extension Commands (Bottom) [Press ':' to focus] "))
	termSb.WriteString("\n")
	if m.LastCmdOut != "" {
		outLines := strings.Split(m.LastCmdOut, "\n")
		if len(outLines) > 2 {
			outLines = outLines[len(outLines)-2:]
		}
		termSb.WriteString(m.Styles.MutedText.Render(fmt.Sprintf("  Output: %s\n", strings.Join(outLines, " | "))))
	}
	termSb.WriteString("  " + m.TermInput.View())

	termBox := termBorderStyle.Width(m.Width - 2).Height(5).Render(termSb.String())

	// --- Bottom Status Bar ---
	statusLeft := " [Tab/Shift+Tab] Focus Pane  •  [/] Fuzzy Filter  •  [Alt+D] Disk Analyzer  •  [:] Command  •  [.] Hidden"
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

	if m.State == StateDialog {
		return m.Dialog.View(m.Styles, m.Width, m.Height)
	}

	return viewStr
}

func (m AppModel) renderAnalyzerView() string {
	var sb strings.Builder
	sb.WriteString(m.Styles.ModalTitle.Render(fmt.Sprintf("📊 Disk Space Analyzer: %s", m.AnalyzerResult.DirPath)))
	sb.WriteString("\n")
	sb.WriteString(m.Styles.MutedText.Render(fmt.Sprintf("   Total Directory Size: %s", analyzer.FormatBytes(m.AnalyzerResult.TotalSize))))
	sb.WriteString("\n\n")

	maxRows := m.Height - 12
	items := m.AnalyzerResult.Items
	endIdx := min(len(items), maxRows)

	for i := 0; i < endIdx; i++ {
		item := items[i]
		prefix := "📁"
		if !item.IsDir {
			prefix = "📄"
		}
		line := fmt.Sprintf("  %s %-20s  %s  %5.1f%%  %-10s", prefix, item.Name, item.Bar, item.Percentage, analyzer.FormatBytes(item.SizeBytes))
		if i == m.AnalyzerCursor {
			sb.WriteString(m.Styles.SelectedItem.Render(fmt.Sprintf("> %s", line)))
		} else {
			sb.WriteString(m.Styles.Item.Render(fmt.Sprintf("  %s", line)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(m.Styles.MutedText.Render(" [↑/↓] Navigate  •  [Enter] Open Dir  •  [Alt+D / Esc] Close Analyzer"))

	box := m.Styles.ModalBox.Render(sb.String())

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars(" "),
	)
}
