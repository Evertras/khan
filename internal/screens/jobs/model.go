package jobs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/datatree"
	"github.com/evertras/khan/internal/components/errview"
	"github.com/evertras/khan/internal/components/logs"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/screens/jobs/list"
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

	inspect         *api.Job
	inspectDataTree datatree.Model

	errorMessage errview.Model

	logView logs.Model

	list tea.Model

	activeState state
}

func NewEmptyModel(size screens.Size) Model {
	return Model{
		errorMessage: errview.NewEmptyModel(),
		size:         size,
		list:         list.New(size),
	}
}

func (m Model) Init() tea.Cmd {
	cancelExistingLogStream()

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
		if msg.String() == "esc" {
			m.activeState = stateList

			cancelExistingLogStream()
		}

	case *api.Job:
		m.activeState = stateInspect
		m.inspect = msg
		m.inspectDataTree = datatree.New(m.inspect)
		m.inspectDataTree, _ = m.inspectDataTree.Update(m.size)

	case errMsg:
		m.activeState = stateError
		m.errorMessage = errview.NewModelWithMessage(msg.Error())

	case screens.Size:
		m.size = msg

	case list.ShowLogs:
		m.activeState = stateLogs
		m.logView, _ = m.logView.Update(m.size)
		cmds = append(cmds, showLogsForJobCmd(msg.JobID))
	}

	m.inspectDataTree, cmd = m.inspectDataTree.Update(msg)
	cmds = append(cmds, cmd)

	switch m.activeState {
	case stateError:
		m.errorMessage, cmd = m.errorMessage.Update(msg)

	case stateLogs:
		m, cmd = m.updateLogView(msg)

	case stateInspect:
		m, cmd = m.updateInspect(msg)

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
		return m.viewLogs()

	case stateInspect:
		return m.viewInspect()

	case stateList:
		return m.list.View()

	default:
		return m.list.View()
	}
}
