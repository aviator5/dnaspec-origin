package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// SuccessStyle is for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	// ErrorStyle is for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	// InfoStyle is for informational messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))

	// SubtleStyle is for secondary text
	SubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	// CodeStyle is for file names and code snippets
	CodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6"))
)
