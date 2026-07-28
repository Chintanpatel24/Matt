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
	asciiLogo := `███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓▄▄▄█████▓
▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓  ██▒ ▓▒
▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▒ ▓██░ ▒░
▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ░ ▓██▓ ░ 
▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░   ▒██▒ ░ 
░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░     ▒ ░░   
░  ░      ░  ▒   ▒▒ ░   ░        ░    
░      ░     ░   ▒    ░        ░      
       ░         ░  ░`

	var sb strings.Builder

	// Logo Header
	sb.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Render(m.Styles.HeaderTitle.Render(asciiLogo)))
	sb.WriteString("\n\n")

	// Version Subtitle
	subtitle := fmt.Sprintf("Matt Black Terminal File Manager %s", version.Version)
	sb.WriteString(lipgloss.NewStyle().Width(52).Align(lipgloss.Center).Render(m.Styles.MutedText.Render(subtitle)))
	sb.WriteString("\n\n")

	// Prompt Title
	promptHeader := m.Styles.Header.Render(" Choose Workspace Session ")
	sb.WriteString(lipgloss.NewStyle().Width(52).Align(lipgloss.Center).Render(promptHeader))
	sb.WriteString("\n\n")

	// Options
	opt0Str := fmt.Sprintf("Open Current Directory (%s)", shortenPath(m.CurrentPath))
	opt1Str := fmt.Sprintf("Open Fresh Session (%s)", shortenPath(m.HomePath))

	var opt0, opt1 string
	if m.Cursor == 0 {
		opt0 = m.Styles.ModalButtonActive.Render(fmt.Sprintf(" ▌ %s ", opt0Str))
		opt1 = m.Styles.MutedText.Render(fmt.Sprintf("   %s", opt1Str))
	} else {
		opt0 = m.Styles.MutedText.Render(fmt.Sprintf("   %s", opt0Str))
		opt1 = m.Styles.ModalButtonActive.Render(fmt.Sprintf(" ▌ %s ", opt1Str))
	}

	sb.WriteString(opt0)
	sb.WriteString("\n\n")
	sb.WriteString(opt1)
	sb.WriteString("\n\n")

	// Footer controls hint
	controls := m.Styles.MutedText.Render("[↑/↓] Navigate  •  [Enter] Confirm  •  [q] Quit")
	sb.WriteString(lipgloss.NewStyle().Width(52).Align(lipgloss.Center).Render(controls))

	box := m.Styles.ModalBox.Render(sb.String())

	if m.Width <= 0 || m.Height <= 0 {
		return box
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars(" "),
	)
}

func shortenPath(path string) string {
	home := config.GetHomeDir()
	if strings.HasPrefix(path, home) {
		path = "~" + path[len(home):]
	}
	if len(path) > 28 {
		return "..." + path[len(path)-25:]
	}
	return path
}
