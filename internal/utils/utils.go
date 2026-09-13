package utils

func SliceToMap[T any, K comparable](s []T, key func(T) K) map[K]T {
	m := make(map[K]T, len(s))

	for _, e := range s {
		m[key(e)] = e
	}

	return m
}

func SliceFromMapValue[T any, R any](items map[string]T, fn func(T) R) []R {
	result := make([]R, 0, len(items))

	for _, item := range items {
		result = append(result, fn(item))
	}

	return result
}