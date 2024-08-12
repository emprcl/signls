package core

type Direction uint8

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func (d Direction) Symbol() string {
	switch d {
	case 0:
		return "↑"
	case 1:
		return "→"
	case 2:
		return "↓"
	case 3:
		return "←"
	default:
		return "↑"
	}
}
