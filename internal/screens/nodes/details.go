package nodes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/khan/internal/keyvalsort"
	"github.com/evertras/khan/internal/styles"
)

func (m Model) updateDetails(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.details = nil
		}
	}

	return m, nil
}

func (m Model) viewDetails() string {
	return m.detailsDataTree.View()
}

func (m Model) viewDetailsOld() string {
	d := m.details

	body := strings.Builder{}

	l := func(k, v string) {
		body.WriteString(styles.Key.Render(k+": ") + v + "\n")
	}

	l("id", d.ID)
	l("name", d.Name)
	l("status", d.Status)
	l("elgbl", d.SchedulingEligibility)
	l("drain", fmt.Sprintf("%v", d.Drain))
	l("dc", d.Datacenter)
	l("addr", d.HTTPAddr)

	l("devices", "")

	for _, device := range d.NodeResources.Devices {
		l("  - id", device.ID())

		count := len(device.Instances)

		if count > 1 {
			l("    count:", fmt.Sprintf("%d", count))
		}
	}

	l("meta", "")

	sortedMetaVals := keyvalsort.SortedStringMapValues(d.Meta)

	for _, kv := range sortedMetaVals {
		l("  "+kv.Key, kv.Val)
	}

	body.WriteString("\n")

	return body.String()
}
