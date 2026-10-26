package tui

type hitOption struct {
	label string
	count int
}

var hitOptions = []hitOption{
	{"Once", 1}, {"Twice", 2}, {"3 times", 3}, {"4 times", 4}, {"5 times", 5},
	// aliases
	{"1 hits", 1}, {"2 hits", 2}, {"3 hits", 3}, {"4 hits", 4}, {"5 hits", 5},
}

var metronomeHitOptions = []hitOption{
	{"Never", 0},
	{"Once", 1}, {"Twice", 2}, {"3 times", 3}, {"4 times", 4}, {"5 times", 5},
	{"6 times", 6}, {"7 times", 7}, {"8 times", 8}, {"9 times", 9}, {"10 times", 10},
}

func hitLabels(opts []hitOption) []string {
	labels := make([]string, len(opts))
	for i, o := range opts {
		labels[i] = o.label
	}
	return labels
}

func hitCount(opts []hitOption, label string) int {
	for _, o := range opts {
		if o.label == label {
			return o.count
		}
	}

	return 0
}
