package slice

// Deduplicate returns a slice of the same comparable type
// with all duplicate values removed.
// Note that the order of elements *will not* be preserved.
func Deduplicate[T comparable](s []T) []T {
	table := map[T]struct{}{}
	for _, v := range s {
		table[v] = struct{}{}
	}

	i := 0
	result := make([]T, len(table))
	for k := range table {
		result[i] = k
		i++
	}
	return result
}
