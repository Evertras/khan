package nodes

import tea "github.com/charmbracelet/bubbletea"

type BackMsg struct{}

func BackCmd() tea.Msg {
	return BackMsg{}
}
