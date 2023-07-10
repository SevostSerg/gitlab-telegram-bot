package extensions

func Abs(value int) int {
	if value >= 0 {
		return value
	}

	return value * (-1)
}
