package sliceutil

func Map[T any, R any](collection []T, mapFunc func(T) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = mapFunc(item)
	}

	return result
}
