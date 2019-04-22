package shared

import (
	"sort"
)

type Deduplicable interface {
	sort.Interface
	Eq(i, j int) bool
	Copy(j, i int)
}

func Deduplicate(in Deduplicable, truncate func(len int)) {
	DeduplicateAny(
		func() { sort.Sort(in) },
		in.Len,
		in.Eq,
		in.Copy,
		truncate,
	)
}

func DeduplicateAny(sort func(), len func() int, eq func(i, j int) bool, copy func(j, i int), truncate func(len int)) {
	sort()
	j := 0
	for i := 1; i < len(); i++ {
		if eq(i, j) {
			continue
		}
		j++
		copy(j, i)
	}
	truncate(j + 1)
}
