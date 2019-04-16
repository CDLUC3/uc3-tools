package misc

import (
	"sort"
)

func SortUniq(strs []string) []string {
	if len(strs) == 0 {
		return strs
	}
	duplicate := make([]string, len(strs))
	copy(duplicate, strs)
	return destructiveSort(duplicate)
}

func SliceContains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func destructiveSort(strs []string) []string {
	sort.Strings(strs)
	j := 0
	for i := 1; i < len(strs); i++ {
		if strs[j] == strs[i] {
			continue
		}
		j++
		strs[j] = strs[i]
	}
	return strs[:j+1]
}
