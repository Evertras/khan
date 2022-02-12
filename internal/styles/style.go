package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	term = termenv.EnvColorProfile()
	dot  = colorFg(" • ", "236")

	// Pre-defined colors
	ColorSubtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#888888"}
	ColorHighlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	ColorSpecial   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	// Pre-defined styles
	Good  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	Error = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Bold  = lipgloss.NewStyle().Bold(true)

	Title    = lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	Subtitle = lipgloss.NewStyle().Foreground(ColorSubtle)
)

func Checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

func Header(title, subtitle string) string {
	return fmt.Sprintf("%s\n%s\n-----------------------------------\n",
		Title.Render(title),
		Subtitle.Render(subtitle),
	)
}
