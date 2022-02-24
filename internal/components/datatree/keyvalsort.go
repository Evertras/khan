package datatree

import (
	"reflect"
	"sort"
)

type keyVal struct {
	key string
	val reflect.Value
}

type keyValList []keyVal

// Interface check
var _ sort.Interface = keyValList{}

func (kv keyValList) Len() int {
	return len(kv)
}

func (kv keyValList) Less(i, j int) bool {
	return kv[i].key < kv[j].key
}

func (kv keyValList) Swap(i, j int) {
	tmp := kv[i]
	kv[i] = kv[j]
	kv[j] = tmp
}
