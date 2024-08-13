package core

type Direction uint8

const (
	NONE Direction = iota
	UP
	RIGHT
	DOWN
	LEFT
)

func DirectionFromString(dir string) Direction {
	switch dir {
	case "up":
		return UP
	case "right":
		return RIGHT
	case "down":
		return DOWN
	case "left":
		return LEFT
	default:
		return UP
	}
}

func (d Direction) Symbol() string {
	switch d {
	case NONE:
		return "•"
	case UP:
		return "↑"
	case RIGHT:
		return "→"
	case DOWN:
		return "↓"
	case LEFT:
		return "←"
	default:
		return "↑"
	}
}

func (d Direction) NextPosition(x, y int) (int, int) {
	switch d {
	case NONE:
		return x, y
	case UP:
		return x, y - 1
	case RIGHT:
		return x + 1, y
	case DOWN:
		return x, y + 1
	case LEFT:
		return x - 1, y
	default:
		return 0, 0
	}
}
