package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PermissionDialog represents an interactive approval modal dialog.
type PermissionDialog struct {
	Title       string
	Message     string
	ActiveIndex int // 0 = Confirm, 1 = Cancel
	IsOpen      bool
	OnConfirm   func()
}

// NewPermissionDialog initializes a modal dialog.
func NewPermissionDialog(title, message string, onConfirm func()) PermissionDialog {
	return PermissionDialog{
		Title:       title,
		Message:     message,
		ActiveIndex: 0,
		IsOpen:      true,
		OnConfirm:   onConfirm,
	}
}

// View renders the modal centered in screen dimensions.
func (d PermissionDialog) View(styles Styles, screenWidth, screenHeight int) string {
	if !d.IsOpen {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(styles.ModalTitle.Render("🔒 Permission Required"))
	sb.WriteString("\n\n")
	sb.WriteString(styles.Item.Render(d.Message))
	sb.WriteString("\n\n")

	btnYes := styles.ModalButton.Render("  Allow  ")
	btnNo := styles.ModalButton.Render("  Deny  ")

	if d.ActiveIndex == 0 {
		btnYes = styles.ModalButtonActive.Render(" [ Allow ] ")
	} else {
		btnNo = styles.ModalButtonActive.Render(" [ Deny ] ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, btnYes, "   ", btnNo)
	sb.WriteString(buttons)

	content := styles.ModalBox.Render(sb.String())

	return lipgloss.Place(
		screenWidth,
		screenHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
}

// QuickConfirmPrompt renders a quick inline confirmation text line for statusbar.
func QuickConfirmPrompt(action string) string {
	return fmt.Sprintf("⚠️ Confirm action: %s? Press [y] to approve, [n] or [Esc] to cancel.", action)
}
