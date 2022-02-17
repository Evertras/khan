package keyvalsort

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyValSortOverManyIterations(t *testing.T) {
	tests := []struct {
		input map[string]string
	}{
		{
			input: map[string]string{},
		},
		{
			input: map[string]string{
				"ok": "thing",
			},
		},
		{
			input: nil,
		},
		{
			input: map[string]string{
				"hmm":      "yes",
				"hmm3":     "yes",
				"hmm45":    "yes",
				"hmm1":     "yes",
				"another":  "sorting",
				"surprise": "yay",
				"何":        "カラオケ",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				sorted := SortedStringMapValues(test.input)

				assert.Len(t, sorted, len(test.input), "Length changed")

				if test.input == nil {
					return
				}

				keys := []string{}

				for _, kv := range sorted {
					keys = append(keys, kv.Key)

					assert.Contains(t, test.input, kv.Key, "Missing key in original map")
				}

				assert.True(t, sort.StringsAreSorted(keys), "Keys not sorted")

				if t.Failed() {
					t.FailNow()
				}
			}
		})
	}
}
