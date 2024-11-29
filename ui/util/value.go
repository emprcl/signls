package util

func Clamp(value, minimum, maximum int) int {
	return max(min(value, maximum), minimum)
}

// Mod handles the modulo operation for negative numbers, ensuring
// the result is always non-negative.
func Mod(a, b int) int {
	return (a%b + b) % b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
