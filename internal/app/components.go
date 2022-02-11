package app

import (
	"fmt"

	"github.com/muesli/termenv"
)

var (
	term = termenv.EnvColorProfile()
	dot  = colorFg(" • ", "236")

	// Pre-defined styles
	sTitle   = makeFgStyle("210")
	sKeyword = makeFgStyle("211")
	sSubtle  = makeFgStyle("241")
)

func checkbox(label string, checked bool) string {
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

func header(title, subtitle string) string {
	return fmt.Sprintf("%s\n%s\n-----------------------------------\n",
		sTitle(title),
		sSubtle(subtitle),
	)
}
