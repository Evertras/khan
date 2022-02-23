package joblist

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/khan/internal/styles"
)

func (m Model) updateInspect(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.inspect = nil
		}
	}
	return m, nil
}

func (m Model) viewInspect() string {
	return m.inspectDataTree.View()
}

func (m Model) viewInspectOld() string {
	i := m.inspect

	body := strings.Builder{}

	activeStyle := styles.Key

	l := func(indent int, k string, v string) {
		body.WriteString(strings.Repeat("  ", indent) + activeStyle.Render(k+": ") + v + "\n")
	}

	ln := func(indent int, k string, v *string) {
		if v == nil {
			return
		}

		l(indent, k, *v)
	}

	ln(0, "id", i.ID)
	ln(0, "name", i.Name)
	l(0, "dc", strings.Join(i.Datacenters, ", "))
	ln(0, "status", i.Status)
	if i.SubmitTime != nil {
		submitTime := time.UnixMicro(*i.SubmitTime / 1000)
		l(0, "submitted", fmt.Sprintf("%v", submitTime))
	}

	colors := []lipgloss.AdaptiveColor{
		{
			Dark:  "#8ff",
			Light: "#044",
		},
		{
			Dark:  "#f8f",
			Light: "#400",
		},
		{
			Dark:  "#88f",
			Light: "#004",
		},
	}

	l(0, "taskgroups", "")
	for _, tg := range i.TaskGroups {
		l(1, *tg.Name, "")

		if tg.Count != nil && *tg.Count != 1 {
			l(2, "count:", fmt.Sprintf("%d", *tg.Count))
		}

		l(2, "tasks", "")

		iColor := 0
		nextColor := func() lipgloss.Style {
			iColor++
			return styles.Key.Copy().Foreground(colors[iColor%len(colors)])
		}

		for _, task := range tg.Tasks {
			oldStyle := activeStyle.Copy()
			activeStyle = nextColor()
			l(3, task.Name, "")

			if task.Kind != "" {
				l(4, "kind", task.Kind)
			}

			if len(task.Artifacts) != 0 {
				l(4, "artifacts", "")

				for _, artifact := range task.Artifacts {
					ln(5, "- source", artifact.GetterSource)
					ln(5, "  dest", artifact.RelativeDest)
					ln(5, "  mode", artifact.GetterMode)
				}
			}

			if len(task.Templates) != 0 {
				l(4, "templates", "")

				for _, template := range task.Templates {
					ln(5, "- dest", template.DestPath)
				}
			}

			if task.Resources != nil && len(task.Resources.Devices) > 0 {
				l(4, "devices", "")

				for _, device := range task.Resources.Devices {
					l(5, "name", device.Name)
					if device.Count != nil && *device.Count != 1 {
						l(5, "count", fmt.Sprintf("%d", *device.Count))
					}
				}
			}

			if task.Lifecycle != nil {
				l(4, "lifecycle", task.Lifecycle.Hook)
			}

			l(4, "driver", task.Driver)

			switch task.Driver {
			case "docker":
				l(5, "image", task.Config["image"].(string))

			case "raw_exec":
				l(5, "command", task.Config["command"].(string))
			}

			activeStyle = oldStyle
		}
	}

	return body.String()
}
