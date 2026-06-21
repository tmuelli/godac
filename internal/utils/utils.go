package utils

func SliceToMap[T any, K comparable](s []T, key func(T) K) map[K]T {
	m := make(map[K]T, len(s))

	for _, e := range s {
		m[key(e)] = e
	}

	return m
}