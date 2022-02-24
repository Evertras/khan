package joblist

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/datatree"
	"github.com/evertras/khan/internal/components/errview"
	"github.com/evertras/khan/internal/components/logs"
	"github.com/evertras/khan/internal/screens"
)

type errMsg error

type Model struct {
	size screens.Size

	jobs []*api.JobListStub

	inspect *api.Job

	inspectDataTree datatree.Model

	table table.Model

	showServices bool
	showBatch    bool

	lastUpdated time.Time

	confirmStopIDs []string

	errorMessage errview.Model

	logView logs.Model
}

func NewEmptyModel(size screens.Size) Model {
	return Model{
		showServices: true,
		showBatch:    true,
		table:        genListTable(),
		errorMessage: errview.NewEmptyModel(),
		size:         size,
		lastUpdated:  time.Now(),
	}
}

const (
	tableKeyID     = "id"
	tableKeyName   = "name"
	tableKeyStatus = "status"
)

func NewModelWithJobs(size screens.Size, jobs []*api.JobListStub) Model {
	m := NewEmptyModel(size)

	m.jobs = jobs

	m.table = genListTable().WithRows(m.generateRows())

	return m
}

func (m Model) Init() tea.Cmd {
	return refreshJobsCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case []*api.JobListStub:
		m.jobs = msg
		m.table = m.table.WithRows(m.generateRows())
		m.lastUpdated = time.Now()

	case *api.Job:
		m.inspect = msg
		m.inspectDataTree = datatree.New(m.inspect)

	case errMsg:
		m.errorMessage = errview.NewModelWithMessage(msg.Error())
	}

	m.inspectDataTree, cmd = m.inspectDataTree.Update(msg)
	cmds = append(cmds, cmd)

	switch {
	case m.errorMessage.Active():
		m.errorMessage, cmd = m.errorMessage.Update(msg)

	case len(m.confirmStopIDs) != 0:
		m, cmd = m.updateConfirmStop(msg)

	case logCancel != nil:
		m, cmd = m.updateLogView(msg)

	case m.inspect != nil:
		m, cmd = m.updateInspect(msg)

	default:
		m, cmd = m.updateMainView(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.errorMessage.Active() {
		return m.errorMessage.View()
	}

	if len(m.confirmStopIDs) != 0 {
		return m.viewConfirmStop()
	}

	if logCancel != nil {
		return m.viewLogs()
	}

	if m.inspect != nil {
		return m.viewInspect()
	}

	return m.viewMain()
}
