package util

import (
	"sort"
)

type IntSlice []int64

func (s IntSlice) Len() int { return len(s) }

func (s IntSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }

func Sort(l []int64) []int64 {
	sort.Sort(IntSlice(l))
	return l
}

type StringSlice []string

func (s StringSlice) Len() int { return len(s) }

func (s StringSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s StringSlice) Less(i, j int) bool { return s[i] < s[j] }

func StringSort(l []string) []string {
	sort.Sort(StringSlice(l))
	return l
}
