package calc

const (
	MAX_UNSIGNED_INT_16 = 65535
	MAX_UNSIGNED_INT_32 = 4294967295
)

func OF16(v float64) float64 {
	value := int(v)

	if value > MAX_UNSIGNED_INT_16 {
		value = value % MAX_UNSIGNED_INT_16
	}

	return float64(value)
}

func OF32(v float64) float64 {
	value := int(v)

	if value > MAX_UNSIGNED_INT_32 {
		value = value % MAX_UNSIGNED_INT_32
	}

	return float64(value)
}