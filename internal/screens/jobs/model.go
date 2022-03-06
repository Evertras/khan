package jobs

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/datatree"
	"github.com/evertras/khan/internal/components/errview"
	"github.com/evertras/khan/internal/components/logs"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/screens/jobs/list"
)

type errMsg error

type Model struct {
	size screens.Size

	inspect *api.Job

	inspectDataTree datatree.Model

	showServices bool
	showBatch    bool

	lastUpdated time.Time

	confirmStopIDs []string

	errorMessage errview.Model

	logView logs.Model

	list tea.Model
}

func NewEmptyModel(size screens.Size) Model {
	return Model{
		showServices: true,
		showBatch:    true,
		errorMessage: errview.NewEmptyModel(),
		size:         size,
		lastUpdated:  time.Now(),
		list:         list.New(size),
	}
}

func (m Model) Init() tea.Cmd {
	if logCancel != nil {
		logCancel()
		logCancel = nil
	}

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
	case *api.Job:
		m.inspect = msg
		m.inspectDataTree = datatree.New(m.inspect)
		m.inspectDataTree, _ = m.inspectDataTree.Update(m.size)

	case errMsg:
		m.errorMessage = errview.NewModelWithMessage(msg.Error())

	case screens.Size:
		m.size = msg
	}

	m.inspectDataTree, cmd = m.inspectDataTree.Update(msg)
	cmds = append(cmds, cmd)

	switch {
	case m.errorMessage.Active():
		m.errorMessage, cmd = m.errorMessage.Update(msg)

	case logCancel != nil:
		m, cmd = m.updateLogView(msg)

	case m.inspect != nil:
		m, cmd = m.updateInspect(msg)

	default:
		m.list, cmd = m.list.Update(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.errorMessage.Active() {
		return m.errorMessage.View()
	}

	if logCancel != nil {
		return m.viewLogs()
	}

	if m.inspect != nil {
		return m.viewInspect()
	}

	return m.list.View()
}
