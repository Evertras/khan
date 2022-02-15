package joblist

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/nomad/api"

	"github.com/evertras/khan/internal/styles"
)

func (m Model) updateMainView(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case []*api.JobListStub:
		m = NewModelWithJobs(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			m.showBatch = !m.showBatch
			m.table = m.table.WithRows(m.generateRows())

		case "e":
			m.showServices = !m.showServices
			m.table = m.table.WithRows(m.generateRows()).Focused(true)

		case "g":
			cmds = append(cmds, garbageCollectCmd)

		case "r":
			cmds = append(cmds, refreshJobsCmd)

		case "s":
			ids := []string{}

			for _, row := range m.table.SelectedRows() {
				if id, exists := row.Data[tableKeyID]; exists {
					ids = append(ids, id.(string))
				}
			}

			if len(ids) > 0 {
				m.confirmStopIDs = ids
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) genHelpBox() string {
	deleteHelp := "(s)top selected"
	if len(m.table.SelectedRows()) == 0 {
		deleteHelp = styleSubtle.Render(deleteHelp)
	}

	helpLines := []string{
		"Space/enter to select\n",
		deleteHelp,
		"(g)arbage collect (clears selections)",
		"(r)efresh jobs (clears selections)",
		"(f)ollow logs of random alloc",
	}

	return styleHelp.Render(strings.Join(helpLines, "\n"))
}

func (m Model) viewMain() string {
	if len(m.jobs) == 0 {
		return ""
	}

	tableView := m.table.View()

	body := strings.Builder{}

	body.WriteString(fmt.Sprintf("Last updated: %s\n", m.lastUpdated.Format("2006-01-02 15:04:05")))

	body.WriteString("Filters: ")
	body.WriteString(styles.Checkbox("Show s(e)rvices", m.showServices))
	body.WriteString("  ")
	body.WriteString(styles.Checkbox("Show (b)atch jobs", m.showBatch))
	body.WriteString("\n")
	body.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		tableView,
		m.genHelpBox(),
	))

	return body.String()
}
