package jobs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/screens/jobs/inspect"
	"github.com/evertras/khan/internal/screens/jobs/list"
	"github.com/evertras/khan/internal/screens/jobs/logs"
)

type Model struct {
	size screens.Size

	activeState state

	subviews map[state]tea.Model
}

func New(size screens.Size) Model {
	return Model{
		size: size,
		subviews: map[state]tea.Model{
			stateList:    list.New(size),
			stateLogs:    logs.New(""),
			stateInspect: inspect.New(nil),
		},
	}
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	for _, subview := range m.subviews {
		cmds = append(cmds, subview.Init())
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.activeState = stateList
		}

	case screens.Size:
		m.size = msg

	case *api.Job:
		cmds = append(cmds, m.toStateInspect(msg))

	case list.ShowLogs:
		cmds = append(cmds, m.toStateLogs(msg.JobID))
	}

	m.subviews[m.activeState], cmd = m.subviews[m.activeState].Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.subviews[m.activeState].View()
}
