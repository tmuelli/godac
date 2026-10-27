package tui

import (
	"strconv"
	"strings"
)

func cloneRows(rows [][]string) [][]string {
	output := make([][]string, len(rows))
	for i, r := range rows {
		output[i] = append([]string{}, r...)
	}
	return output
}

func cloneColumns(columns []gridColumn) []gridColumn {
	return append([]gridColumn{}, columns...)
}

func atoiOrZero(value string) int {
	n, _ := strconv.Atoi(value)
	return n
}

func healthPercent(l list) int {
	return atoiOrZero(strings.TrimSuffix(l.values[HEALTH_PERCENT], "%"))
}
