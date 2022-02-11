package styles

import (
	"fmt"

	"github.com/muesli/termenv"
)

var (
	term = termenv.EnvColorProfile()
	dot  = colorFg(" • ", "236")

	// Pre-defined styles
	Title   = makeFgStyle("210")
	Keyword = makeFgStyle("211")
	Subtle  = makeFgStyle("241")
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
		Title(title),
		Subtle(subtitle),
	)
}
