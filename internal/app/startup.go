package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Chintanpatel24/Matt/internal/config"
	"github.com/Chintanpatel24/Matt/internal/ui"
	"github.com/Chintanpatel24/Matt/internal/version"
)

// StartupChoice indicates which directory option the user selected.
type StartupChoice int

const (
	ChoiceCurrentDir StartupChoice = iota
	ChoiceHomeDir
)

// StartupModel renders the initial Matt welcome prompt.
type StartupModel struct {
	Cursor      int
	Selected    StartupChoice
	IsDone      bool
	Width       int
	Height      int
	Config      config.Config
	Styles      ui.Styles
	CurrentPath string
	HomePath    string
}

// NewStartupModel initializes the startup prompt screen.
func NewStartupModel(cfg config.Config, currentDir string) StartupModel {
	return StartupModel{
		Cursor:      0,
		Selected:    ChoiceCurrentDir,
		IsDone:      false,
		Config:      cfg,
		Styles:      ui.NewStyles(cfg),
		CurrentPath: currentDir,
		HomePath:    config.GetHomeDir(),
	}
}

func (m StartupModel) Init() tea.Cmd {
	return nil
}

func (m StartupModel) Update(msg tea.Msg) (StartupModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < 1 {
				m.Cursor++
			}
		case "1":
			m.Cursor = 0
			m.Selected = ChoiceCurrentDir
			m.IsDone = true
		case "2":
			m.Cursor = 1
			m.Selected = ChoiceHomeDir
			m.IsDone = true
		case "enter":
			if m.Cursor == 0 {
				m.Selected = ChoiceCurrentDir
			} else {
				m.Selected = ChoiceHomeDir
			}
			m.IsDone = true
		}
	}
	return m, nil
}

func (m StartupModel) View() string {
	asciiLogo := `
  ███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓▄▄▄█████▓
 ▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓  ██▒ ▓▒
 ▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▒ ▓██░ ▒░
 ▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ░ ▓██▓ ░ 
 ▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░   ▒██▒ ░ 
 ░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░     ▒ ░░   
`

	var sb strings.Builder

	// Logo Header Box
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.Config.Theme.Accent)).
		Bold(true).
		Align(lipgloss.Center)

	sb.WriteString(logoStyle.Render(asciiLogo))
	sb.WriteString("\n")

	// Version Badge
	verBadge := lipgloss.NewStyle().
		Background(lipgloss.Color(m.Config.Theme.Accent)).
		Foreground(lipgloss.Color(m.Config.Theme.Bg)).
		Bold(true).
		Padding(0, 2).
		Render(" MATT " + version.Version + " ")

	subTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.Config.Theme.TextMuted)).
		Render("High-Performance Terminal Workspace & File Manager")

	sb.WriteString(lipgloss.NewStyle().Width(56).Align(lipgloss.Center).Render(verBadge))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Width(56).Align(lipgloss.Center).Render(subTitle))
	sb.WriteString("\n\n")

	// Prompt Title
	headerText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.Config.Theme.TextPrimary)).
		Bold(true).
		Render("Select Workspace Session:")

	sb.WriteString("  " + headerText + "\n\n")

	// Session Options
	opt0Path := shortenPath(m.CurrentPath)
	opt1Path := shortenPath(m.HomePath)

	activeCardStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.Config.Theme.Accent)).
		Background(lipgloss.Color(m.Config.Theme.Selection)).
		Foreground(lipgloss.Color(m.Config.Theme.TextPrimary)).
		Padding(0, 1).
		Width(52)

	inactiveCardStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(m.Config.Theme.Border)).
		Background(lipgloss.Color(m.Config.Theme.Bg)).
		Foreground(lipgloss.Color(m.Config.Theme.TextMuted)).
		Padding(0, 1).
		Width(52)

	var card0, card1 string
	if m.Cursor == 0 {
		card0 = activeCardStyle.Render(fmt.Sprintf("➜ [1] Open Current Directory\n    %s", opt0Path))
		card1 = inactiveCardStyle.Render(fmt.Sprintf("  [2] Open Home Directory\n    %s", opt1Path))
	} else {
		card0 = inactiveCardStyle.Render(fmt.Sprintf("  [1] Open Current Directory\n    %s", opt0Path))
		card1 = activeCardStyle.Render(fmt.Sprintf("➜ [2] Open Home Directory\n    %s", opt1Path))
	}

	sb.WriteString("  " + strings.ReplaceAll(card0, "\n", "\n  ") + "\n")
	sb.WriteString("  " + strings.ReplaceAll(card1, "\n", "\n  ") + "\n\n")

	// Quick Key Guide Footer
	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.Config.Theme.TextMuted)).
		Render(" [1/2] Quick Select • [↑/↓/k/j] Navigate • [Enter] Launch ")

	sb.WriteString(lipgloss.NewStyle().Width(56).Align(lipgloss.Center).Render(controls))

	mainBox := m.Styles.ModalBox.Width(60).Render(sb.String())

	if m.Width <= 0 || m.Height <= 0 {
		return mainBox
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		mainBox,
		lipgloss.WithWhitespaceChars(" "),
	)
}

func shortenPath(path string) string {
	home := config.GetHomeDir()
	if strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}
	if len(path) > 38 {
		return "..." + path[len(path)-35:]
	}
	return path
}
