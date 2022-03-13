package jobs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/screens/jobs/inspect"
	"github.com/evertras/khan/internal/screens/jobs/logs"
)

type state int

const (
	stateList state = iota
	stateLogs
	stateInspect
)

func (m *Model) refreshSizeCmd() tea.Cmd {
	return func() tea.Msg {
		return m.size
	}
}

func (m *Model) toStateLogs(jobID string) tea.Cmd {
	m.activeState = stateLogs
	m.subviews[stateLogs] = logs.New(jobID)

	return tea.Batch(
		m.subviews[stateLogs].Init(),
		m.refreshSizeCmd(),
	)
}

func (m *Model) toStateInspect(job *api.Job) tea.Cmd {
	m.activeState = stateInspect
	m.subviews[stateInspect] = inspect.New(job, m.size)

	return tea.Batch(
		m.subviews[stateInspect].Init(),
		m.refreshSizeCmd(),
	)
}
