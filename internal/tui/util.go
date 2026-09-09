package tui

func cloneRows(rows [][]string) [][]string{
	output := make([][]string, len(rows))
	for i, r := range rows {
		output[i] = append([]string{}, r...)
	}
	return output
}