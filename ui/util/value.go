package util

func Clamp(value, minimum, maximum int) int {
	return max(min(value, maximum), minimum)
}
