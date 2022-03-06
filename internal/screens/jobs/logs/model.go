package logs

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertras/khan/internal/components/screenviewport"
	"github.com/evertras/khan/internal/styles"
)

type Model struct {
	jobID     string
	allocID   string
	taskGroup string
	task      string

	viewport screenviewport.Model
}

func New(jobID string) Model {
	return Model{
		jobID: jobID,
	}
}

func (m Model) Init() tea.Cmd {
	if m.jobID == "" {
		return nil
	}

	return showLogsForJobCmd(m.jobID)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case logStreamingMsg:
		m.allocID = msg.allocID
		m.taskGroup = msg.taskGroup
		m.task = msg.task

		header := fmt.Sprintf(
			"Logs - %s %s %s/%s\n",
			styles.Header.Render(m.jobID),
			styles.Subtitle.Render(m.allocID),
			styles.Error.Render(m.taskGroup),
			styles.Good.Render(m.task),
		)

		m.viewport = m.viewport.WithHeader(header)
		cmds = append(cmds, func() tea.Msg {
			return logReceivedMsg{
				logStreamingMsg:  msg,
				receivedLogChunk: "",
			}
		})

	case logReceivedMsg:
		m.viewport = m.viewport.Append(msg.receivedLogChunk)
		cmds = append(cmds, func() tea.Msg {
			select {
			case frame, ok := <-msg.logs:
				if !ok {
					return nil
				}

				return logReceivedMsg{
					logStreamingMsg: msg.logStreamingMsg,

					receivedLogChunk: string(frame.Data),
				}

			case err, ok := <-msg.errs:
				if !ok {
					return nil
				}

				return err
			}
		})
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}
