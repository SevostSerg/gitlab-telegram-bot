package extensions

// Returns "enabled" if val = 1 else "disabled"
func GetOptionStatusString(val int) string {
	if val == 1 {
		return "enabled"
	}

	return "disabled"
}
