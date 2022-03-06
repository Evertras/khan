package jobs

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) refreshSizeCmd() tea.Cmd {
	return func() tea.Msg {
		return m.size
	}
}
