package joblist

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/components/errview"
	"github.com/evertras/khan/internal/components/logs"
	"github.com/evertras/khan/internal/screens"
)

type errMsg error

type Model struct {
	size screens.Size

	jobs []*api.JobListStub

	inspect *api.Job

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
	switch msg := msg.(type) {
	case []*api.JobListStub:
		m.jobs = msg
		m.table = m.table.WithRows(m.generateRows())
		m.lastUpdated = time.Now()

	case *api.Job:
		m.inspect = msg

	case errMsg:
		m.errorMessage = errview.NewModelWithMessage(msg.Error())
	}

	if m.errorMessage.Active() {
		var cmd tea.Cmd
		m.errorMessage, cmd = m.errorMessage.Update(msg)
		return m, cmd
	}

	if len(m.confirmStopIDs) != 0 {
		return m.updateConfirmStop(msg)
	}

	if logCancel != nil {
		return m.updateLogView(msg)
	}

	if m.inspect != nil {
		return m.updateInspect(msg)
	}

	return m.updateMainView(msg)
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
