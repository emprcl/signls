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

func (d Direction) NextPosition(x, y int) (int, int) {
	switch d {
	case 0:
		return x, y - 1
	case 1:
		return x + 1, y
	case 2:
		return x, y + 1
	case 3:
		return x - 1, y
	default:
		return 0, 0
	}
}
