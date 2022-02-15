package logs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/screens"
)

type Model struct {
	viewport viewport.Model
	ready    bool

	jobID     string
	allocID   string
	taskGroup string
	task      string

	contents string

	useHighPerformanceRenderer bool
}

func NewJobLogs(jobID string) Model {
	return Model{
		jobID:    jobID,
		contents: "Loading...",

		useHighPerformanceRenderer: true,
	}
}

func (m Model) WithJobInfo(jobID, allocID, taskGroup, task string) Model {
	m.jobID = jobID
	m.allocID = allocID
	m.taskGroup = taskGroup
	m.task = task
	return m
}

func (m Model) Append(data string) Model {
	m.contents += data
	m.viewport.SetContent(m.contents)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.viewport.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case screens.Size:
		headerHeight := 2
		footerHeight := lipgloss.Height("x")
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.contents)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if m.useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	if m.jobID == "" {
		body.WriteString("Logs loading...")
	} else {
		body.WriteString(fmt.Sprintf("%s - %s %s/%s\n", m.jobID, m.allocID, m.taskGroup, m.task))
	}
	body.WriteString(strings.Repeat("=", m.viewport.Width))
	body.WriteString("\n")
	body.WriteString(m.viewport.View())
	return body.String()
}
