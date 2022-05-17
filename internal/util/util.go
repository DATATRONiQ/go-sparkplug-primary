package util

// Returns true iff any element of the first parameter is equal to the second parameter
func Contains[T comparable](items []T, item T) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
