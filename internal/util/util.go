package util

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Returns true iff any element of the first parameter is equal to the second parameter
func Contains[T comparable](items []T, item T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

// Returns a sorted slice of the keys of the given map
func SortedKeys[K constraints.Ordered, V any](m map[K]V) *[]K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return &keys
}

// Maps a function over the given slice and returns a new slice with the results
func MapSlice[T, R any](items *[]T, mapFunc func(T) R) *[]R {
	if items == nil {
		return nil
	}
	result := make([]R, 0, len(*items))
	for _, item := range *items {
		result = append(result, mapFunc(item))
	}
	return &result
}
