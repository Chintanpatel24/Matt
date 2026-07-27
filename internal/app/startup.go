package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Chintanpatel24/Matt/internal/config"
	"github.com/Chintanpatel24/Matt/internal/ui"
)

// StartupChoice indicates which directory option the user selected.
type StartupChoice int

const (
	ChoiceCurrentDir StartupChoice = iota
	ChoiceHomeDir
)

// StartupModel renders the initial Matt welcome prompt.
type StartupModel struct {
	Cursor     int
	Selected   StartupChoice
	IsDone     bool
	Width      int
	Height     int
	Config     config.Config
	Styles     ui.Styles
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
	asciiLogo := `
  ███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓▄▄▄█████▓
  ▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓  ██▒ ▓▒
  ▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▒ ▓██░ ▒░
  ▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ░ ▓██▓ ░ 
  ▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░   ▒██▒ ░ 
  ░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░     ▒ ░░   
  ░  ░      ░  ▒   ▒▒ ░   ░        ░    
  ░      ░     ░   ▒    ░        ░      
         ░         ░  ░`

	var sb strings.Builder
	sb.WriteString(m.Styles.HeaderTitle.Render(asciiLogo))
	sb.WriteString("\n\n")
	sb.WriteString(m.Styles.MutedText.Render("       Matt Black Terminal File Manager v1.0.0"))
	sb.WriteString("\n")
	sb.WriteString(m.Styles.MutedText.Render("     ============================================"))
	sb.WriteString("\n\n")
	sb.WriteString(m.Styles.Header.Render(" Choose your starting workspace session:"))
	sb.WriteString("\n\n")

	opt0 := fmt.Sprintf("   Open Current Directory   (%s)", m.CurrentPath)
	opt1 := fmt.Sprintf("   Open Fresh Session       (%s)", m.HomePath)

	if m.Cursor == 0 {
		opt0 = m.Styles.ModalButtonActive.Render(fmt.Sprintf(" > Open Current Directory   (%s) ", m.CurrentPath))
		opt1 = m.Styles.MutedText.Render(fmt.Sprintf("   Open Fresh Session       (%s)", m.HomePath))
	} else {
		opt0 = m.Styles.MutedText.Render(fmt.Sprintf("   Open Current Directory   (%s)", m.CurrentPath))
		opt1 = m.Styles.ModalButtonActive.Render(fmt.Sprintf(" > Open Fresh Session       (%s) ", m.HomePath))
	}

	sb.WriteString(opt0)
	sb.WriteString("\n\n")
	sb.WriteString(opt1)
	sb.WriteString("\n\n")
	sb.WriteString(m.Styles.MutedText.Render(" [↑/↓] Navigate  •  [Enter] Confirm Choice  •  [q/Ctrl+C] Exit"))

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
