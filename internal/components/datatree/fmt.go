package datatree

import "strings"

func trimNewline(s string) string {
	return strings.TrimSuffix(s, "\n")
}
