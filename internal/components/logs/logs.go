package logs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/styles"
	"github.com/muesli/reflow/wordwrap"
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
	m.contents = ""
	return m
}

func (m Model) Append(data string) Model {
	m.contents += data
	m.viewport.SetContent(wordwrap.String(m.contents, m.viewport.Width))

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

	case screens.Size:
		headerHeight := 2
		footerHeight := 1
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
			m.viewport.SetContent(wordwrap.String(m.contents, msg.Width))
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.viewport.SetContent(wordwrap.String(m.contents, msg.Width))
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) footerView() string {
	footerText := fmt.Sprintf("<%3.f%%>", m.viewport.ScrollPercent()*100)

	info := styles.Subtitle.Render(footerText)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, info, line)
}

func (m Model) View() string {
	body := strings.Builder{}

	if m.jobID == "" {
		body.WriteString("Logs loading...")
	} else {
		jobRow := fmt.Sprintf(
			"Logs - %s %s %s/%s\n",
			styles.Header.Render(m.jobID),
			styles.Subtitle.Render(m.allocID),
			styles.Error.Render(m.taskGroup),
			styles.Good.Render(m.task),
		)
		jobRow += styles.Title.Render(strings.Repeat("─", m.viewport.Width))
		body.WriteString(jobRow)
	}

	body.WriteString("\n")
	body.WriteString(m.viewport.View())
	body.WriteString("\n")
	body.WriteString(m.footerView())

	return body.String()
}
