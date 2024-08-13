package ui

func clamp(value, minimum, maximum int) int {
	return max(min(value, maximum), minimum)
}
