package nodes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/screens/nodes/details"
	"github.com/evertras/khan/internal/screens/nodes/list"
)

type errMsg error

type Model struct {
	size screens.Size

	activeView tea.Model
}

func New(size screens.Size) Model {
	return Model{
		size: size,

		activeView: list.New(size),
	}
}

func (m Model) Init() tea.Cmd {
	return m.activeView.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case list.DetailsSelectedID:
		m.activeView = details.New(string(msg), m.size)

		cmd = m.activeView.Init()
		cmds = append(cmds, cmd)

	case details.Exit:
		m.activeView = list.New(m.size)

		cmd = m.activeView.Init()
		cmds = append(cmds, cmd)

	case screens.Size:
		m.size = msg
	}

	m.activeView, cmd = m.activeView.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(m.activeView.View())

	return body.String()
}
