package datatree

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"github.com/evertras/khan/internal/screens"
)

var (
	styleFieldKey = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#224",
		Dark:  "#b8e",
	}).Bold(true).MarginRight(1)
)

type Model struct {
	data       interface{}
	indentSize int
	showZero   bool

	viewport viewport.Model
	contents string
	ready    bool
}

func New(data interface{}) Model {
	model := Model{
		data:       data,
		indentSize: 2,
	}

	model.updateContents()

	return model
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
		// Some additional navigation keys on top of the default
		switch msg.String() {
		case "g":
			m.viewport.GotoTop()

		case "G":
			m.viewport.GotoBottom()
		}

	case screens.Size:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.KeyMap.PageDown.SetKeys("f", " ", "ctrl+f", "pgdown")
			m.viewport.KeyMap.PageUp.SetKeys("b", "ctrl+b", "pgup")
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(wordwrap.String(m.contents, msg.Width))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
			m.viewport.SetContent(wordwrap.String(m.contents, msg.Width))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(m.viewport.View())

	return body.String()
}
