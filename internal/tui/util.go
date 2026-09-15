package tui

import "strconv"

func cloneRows(rows [][]string) [][]string{
	output := make([][]string, len(rows))
	for i, r := range rows {
		output[i] = append([]string{}, r...)
	}
	return output
}

func atoiOrZero(value string) int {
	n, _ := strconv.Atoi(value)
	return n
}
