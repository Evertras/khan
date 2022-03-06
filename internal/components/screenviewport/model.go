package screenviewport

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"github.com/evertras/khan/internal/screens"
	"github.com/evertras/khan/internal/styles"
)

type Model struct {
	viewport viewport.Model
	ready    bool

	header   string
	contents string

	useHighPerformanceRenderer bool
}

func New(jobID string) Model {
	m := Model{
		useHighPerformanceRenderer: true,
	}

	return m
}

func (m Model) WithHeader(header string) Model {
	m.header = header

	return m
}

func (m Model) SetContent(content string) Model {
	m.viewport.SetContent(wordwrap.String(m.contents, m.viewport.Width))

	return m
}

func (m Model) Append(content string) Model {
	m.contents += content
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
		// Some additional navigation keys on top of the default
		switch msg.String() {
		case "g":
			m.viewport.GotoTop()

		case "G":
			m.viewport.GotoBottom()
		}

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
			m.viewport.KeyMap.PageDown.SetKeys("f", " ", "ctrl+f", "pgdown")
			m.viewport.KeyMap.PageUp.SetKeys("b", "ctrl+b", "pgup")
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

func (m Model) footerView() string {
	footerText := fmt.Sprintf("[%3.f%%]", m.viewport.ScrollPercent()*100)

	const offset = 1
	prefixLine := strings.Repeat("─", offset)
	info := styles.Subtitle.Render(footerText)
	suffixLine := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info))-offset)

	return lipgloss.JoinHorizontal(lipgloss.Center, prefixLine, info, suffixLine)
}

func (m Model) View() string {
	if !m.ready {
		return "Waiting for size..."
	}

	body := strings.Builder{}

	body.WriteString(m.header)
	body.WriteString(styles.Title.Render(strings.Repeat("─", m.viewport.Width)))
	body.WriteString(m.viewport.View())
	body.WriteString("\n")
	body.WriteString(m.footerView())

	return body.String()
}
