package joblist

import (
	"context"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/khan/internal/components/errview"
)

var (
	logCtx    context.Context
	logCancel context.CancelFunc
	logMu     sync.Mutex
)

type logReceivedMsg struct {
	logStreamingMsg

	receivedLogChunk string
}

type logErrMsg error

func (m Model) updateLogView(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case logStreamingMsg:
		m.logView = m.logView.WithJobInfo(msg.jobID, msg.allocID, msg.taskGroup, msg.task)
		cmds = append(cmds, func() tea.Msg {
			select {
			case frame, ok := <-msg.logs:
				if !ok {
					return nil
				}

				return logReceivedMsg{
					logStreamingMsg: msg,

					receivedLogChunk: string(frame.Data),
				}

			case err, ok := <-msg.errs:
				if !ok {
					return nil
				}

				return logErrMsg(err)
			}
		})

	case logReceivedMsg:
		m.logView = m.logView.Append(msg.receivedLogChunk)

	case logErrMsg:
		m.errorMessage = errview.NewModelWithMessage(msg.Error())

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			logMu.Lock()
			if logCancel != nil {
				logCancel()
				logCancel = nil
			}
			logMu.Unlock()
		}
	}

	m.logView, cmd = m.logView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) viewLogs() string {
	return m.logView.View()
}
