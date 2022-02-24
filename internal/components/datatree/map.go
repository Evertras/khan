package datatree

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func (m Model) renderDataNodeMap(data reflect.Value, indentLevel int) string {
	result := strings.Builder{}
	indent := strings.Repeat(" ", indentLevel*m.indentSize)

	result.WriteString(indent)

	iter := data.MapRange()

	keyVals := keyValList{}

	for iter.Next() {
		keyVals = append(keyVals, keyVal{
			key: fmt.Sprintf("%v", iter.Key()),
			val: iter.Value(),
		})
	}

	sort.Sort(keyVals)

	for _, kv := range keyVals {
		result.WriteString("\n")
		keyStr := styleFieldKey.Render(kv.key + ":")
		result.WriteString(indent + keyStr)
		result.WriteString(m.renderDataNode(kv.val, indentLevel+1))
	}

	return trimNewline(result.String())
}
