package jobs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/errview"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/screens/jobs/inspect"
	"github.com/evertras/khan/internal/screens/jobs/list"
	"github.com/evertras/khan/internal/screens/jobs/logs"
)

type state int

const (
	stateList state = iota
	stateLogs
	stateInspect
	stateError
)

type errMsg error

type Model struct {
	size screens.Size

	activeState state

	errorMessage errview.Model

	list    tea.Model
	inspect tea.Model
	logs    tea.Model
}

func New(size screens.Size) Model {
	return Model{
		errorMessage: errview.NewEmptyModel(),
		size:         size,
		list:         list.New(size),
	}
}

func (m Model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	cmds = append(cmds, m.list.Init())

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

	case *api.Job:
		m.activeState = stateInspect
		m.inspect = inspect.New(msg)
		cmds = append(cmds, m.refreshSizeCmd())

	case errMsg:
		m.activeState = stateError
		m.errorMessage = errview.NewModelWithMessage(msg.Error())

	case screens.Size:
		m.size = msg

	case list.ShowLogs:
		m.activeState = stateLogs
		m.logs = logs.New(msg.JobID)
		cmds = append(cmds, m.logs.Init())
		cmds = append(cmds, m.refreshSizeCmd())
	}

	switch m.activeState {
	case stateError:
		m.errorMessage, cmd = m.errorMessage.Update(msg)

	case stateLogs:
		m.logs, cmd = m.logs.Update(msg)

	case stateInspect:
		m.inspect, cmd = m.inspect.Update(msg)

	default:
		m.list, cmd = m.list.Update(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.activeState {
	case stateError:
		return m.errorMessage.View()

	case stateLogs:
		return m.logs.View()

	case stateInspect:
		return m.inspect.View()

	case stateList:
		return m.list.View()

	default:
		return m.list.View()
	}
}
