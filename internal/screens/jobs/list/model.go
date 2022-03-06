package list

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/evertras/khan/internal/components/errview"
	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/styles"
	"github.com/hashicorp/nomad/api"
)

const (
	tableKeyID         = "id"
	tableKeyName       = "name"
	tableKeyStatus     = "status"
	tableKeySubmitDate = "submitDate"
)

var (
	styleHelp   = lipgloss.NewStyle().Width(70).Padding(1).Foreground(styles.ColorSpecial)
	styleSubtle = lipgloss.NewStyle().Foreground(styles.ColorSubtle)

	styleConfirmWarning = styles.Error.Copy().Padding(2)
)

type Model struct {
	size           screens.Size
	jobs           []*api.JobListStub
	table          table.Model
	lastUpdated    time.Time
	errorMessage   errview.Model
	showBatch      bool
	showServices   bool
	confirmStopIDs []string
}

func New(size screens.Size) Model {
	return Model{
		size:         size,
		table:        genListTable(),
		showBatch:    true,
		showServices: true,
	}
}

func (m Model) Init() tea.Cmd {
	return refreshJobsCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if len(m.confirmStopIDs) != 0 {
		return m.updateConfirmStop(msg)
	}

	switch msg := msg.(type) {
	case screens.Size:
		m.size = msg

	case []*api.JobListStub:
		m.jobs = msg
		m.table = m.table.WithRows(m.generateRows())
		m.lastUpdated = time.Now()

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

		case "i":
			if len(m.jobs) == 0 {
				break
			}

			cmds = append(cmds, inspectJobCmd(m.table.HighlightedRow().Data[tableKeyID].(string)))

		case "f":
			if len(m.jobs) == 0 {
				break
			}

			jobID := m.table.HighlightedRow().Data[tableKeyID].(string)
			cmds = append(cmds, func() tea.Msg {
				return ShowLogs{jobID}
			})

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

func (m Model) View() string {
	if len(m.confirmStopIDs) != 0 {
		return m.viewConfirmStop()
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

func (m Model) genHelpBox() string {
	deleteHelp := "(s)top selected"
	if len(m.table.SelectedRows()) == 0 {
		deleteHelp = styleSubtle.Render(deleteHelp)
	}

	helpLines := []string{
		"Space/enter to select\n",
		deleteHelp,
		"(i)nspect job",
		"(g)arbage collect (clears selections)",
		"(r)efresh jobs (clears selections)",
		"(f)ollow logs of random alloc",
	}

	return styleHelp.Render(strings.Join(helpLines, "\n"))
}
