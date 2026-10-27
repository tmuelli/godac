package utils

import "math"

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

// shared because is needed in calc and tui package
func GetMaxHP(base, iv, ev float64, level int) int {
	return int(math.Floor(((2*base)+iv+math.Floor(ev/4))*(float64(level)/100)) + float64(level) + 10)
}
