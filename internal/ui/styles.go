package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/Chintanpatel24/Matt/internal/config"
)

// Styles holds Lip Gloss style definitions for Matt TUI.
type Styles struct {
	Header            lipgloss.Style
	HeaderTitle       lipgloss.Style
	HeaderPath        lipgloss.Style
	HeaderBadge       lipgloss.Style
	PaneTitleActive   lipgloss.Style
	PaneTitleInactive lipgloss.Style
	ActiveBorder      lipgloss.Style
	InactiveBorder    lipgloss.Style
	Item              lipgloss.Style
	SelectedItem      lipgloss.Style
	DirectoryItem     lipgloss.Style
	ExecutableItem    lipgloss.Style
	MutedText         lipgloss.Style
	StatusBar         lipgloss.Style
	PromptLabel       lipgloss.Style
	ModalBox          lipgloss.Style
	ModalTitle        lipgloss.Style
	ModalButton       lipgloss.Style
	ModalButtonActive lipgloss.Style
	MetaLabel         lipgloss.Style
	MetaValue         lipgloss.Style
}

// NewStyles constructs all styles using the Lite Grey & Matt Black theme configuration.
func NewStyles(cfg config.Config) Styles {
	t := cfg.Theme

	return Styles{
		Header: lipgloss.NewStyle().
			Background(lipgloss.Color(t.BgSurface)).
			Foreground(lipgloss.Color(t.TextPrimary)).
			Padding(0, 1).
			Bold(true),

		HeaderTitle: lipgloss.NewStyle().
			Background(lipgloss.Color(t.BgSurface)).
			Foreground(lipgloss.Color(t.Accent)).
			Bold(true),

		HeaderPath: lipgloss.NewStyle().
			Background(lipgloss.Color(t.BgSurface)).
			Foreground(lipgloss.Color(t.Directory)).
			Bold(false),

		HeaderBadge: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Accent)).
			Foreground(lipgloss.Color(t.Bg)).
			Padding(0, 1).
			Bold(true),

		PaneTitleActive: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Accent)).
			Foreground(lipgloss.Color(t.Bg)).
			Padding(0, 1).
			Bold(true),

		PaneTitleInactive: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Selection)).
			Foreground(lipgloss.Color(t.TextMuted)).
			Padding(0, 1),

		ActiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.BorderActive)).
			BorderBackground(lipgloss.Color(t.Bg)).
			Background(lipgloss.Color(t.Bg)),

		InactiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.Border)).
			BorderBackground(lipgloss.Color(t.Bg)).
			Background(lipgloss.Color(t.Bg)),

		Item: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.TextPrimary)).
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Selection)).
			Foreground(lipgloss.Color(t.TextPrimary)).
			Bold(true).
			Padding(0, 1),

		DirectoryItem: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.Directory)).
			Bold(true),

		ExecutableItem: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.Executable)).
			Bold(true),

		MutedText: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.TextMuted)),

		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color(t.BgSurface)).
			Foreground(lipgloss.Color(t.TextMuted)).
			Padding(0, 1),

		PromptLabel: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.Accent)).
			Bold(true),

		ModalBox: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color(t.Accent)).
			BorderBackground(lipgloss.Color(t.BgSurface)).
			Background(lipgloss.Color(t.BgSurface)).
			Padding(1, 2),

		ModalTitle: lipgloss.NewStyle().
			Background(lipgloss.Color(t.BgSurface)).
			Foreground(lipgloss.Color(t.Accent)).
			Bold(true),

		ModalButton: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Selection)).
			Foreground(lipgloss.Color(t.TextMuted)).
			Padding(0, 2),

		ModalButtonActive: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Accent)).
			Foreground(lipgloss.Color(t.Bg)).
			Bold(true).
			Padding(0, 2),

		MetaLabel: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.TextMuted)).
			Bold(true),

		MetaValue: lipgloss.NewStyle().
			Background(lipgloss.Color(t.Bg)).
			Foreground(lipgloss.Color(t.TextPrimary)),
	}
}
