package datatree

import (
	"fmt"
	"reflect"
)

func (m Model) renderDataNode(data reflect.Value, indentLevel int) string {
	for data.Kind() == reflect.Ptr {
		if data.IsNil() {
			return "<nil>"
		}

		data = data.Elem()
	}

	var result string

	switch data.Kind() {
	case reflect.Struct:
		result = m.renderDataNodeStruct(data, indentLevel)

	case reflect.Map:
		result = m.renderDataNodeMap(data, indentLevel)

	case reflect.Array, reflect.Slice:
		result = m.renderDataNodeArray(data, indentLevel)

	default:
		result = fmt.Sprintf("%v", data)
	}

	return trimNewline(result)
}
