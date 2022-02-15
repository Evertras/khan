package joblist

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/styles"
)

var (
	styleHelp   = lipgloss.NewStyle().Width(70).Padding(1).Foreground(styles.ColorSpecial)
	styleSubtle = lipgloss.NewStyle().Foreground(styles.ColorSubtle)

	styleConfirmWarning = styles.Error.Copy().Padding(2)
)
