package home

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	body := strings.Builder{}
	body.WriteString(" Welcome to Khan!  Press the first letter (with shift) of the tabs above to visit each tab.\n\n")
	body.WriteString(" Press 'q' or ctrl+C at any time to quit.")

	body.WriteString("\n")

	return body.String()
}
