package keyvalsort

import "sort"

type KeyVal struct {
	Key string
	Val string
}

type KeyValList []KeyVal

// Interface check
var _ sort.Interface = KeyValList{}

func (kv KeyValList) Len() int {
	return len(kv)
}

func (kv KeyValList) Less(i, j int) bool {
	return kv[i].Key < kv[j].Key
}

func (kv KeyValList) Swap(i, j int) {
	tmp := kv[i]
	kv[i] = kv[j]
	kv[j] = tmp
}

func SortedStringMapValues(m map[string]string) KeyValList {
	vals := []KeyVal{}

	for k, v := range m {
		vals = append(vals, KeyVal{
			Key: k,
			Val: v,
		})
	}

	sort.Sort(KeyValList(vals))

	return vals
}
