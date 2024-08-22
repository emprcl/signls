package core

import "math/bits"

type Direction int

const (
	NONE Direction = 0
	UP   Direction = 1 << iota
	RIGHT
	DOWN
	LEFT
)

var (
	allDirections = []Direction{UP, RIGHT, LEFT, DOWN}
	symbols       = map[Direction]string{
		NONE:                     "•",
		UP:                       "╹",
		DOWN:                     "╻",
		LEFT:                     "╸",
		RIGHT:                    "╺",
		UP | LEFT:                "┛",
		UP | RIGHT:               "┗",
		DOWN | LEFT:              "┓",
		DOWN | RIGHT:             "┏",
		UP | DOWN:                "┃",
		LEFT | RIGHT:             "━",
		UP | LEFT | RIGHT:        "┻",
		DOWN | LEFT | RIGHT:      "┳",
		UP | DOWN | LEFT:         "┫",
		UP | DOWN | RIGHT:        "┣",
		UP | DOWN | LEFT | RIGHT: "╋",
	}
)

func (d Direction) Decompose() []Direction {
	directions := []Direction{}
	for _, dir := range allDirections {
		if !d.Contains(dir) {
			continue
		}
		directions = append(directions, dir)
	}
	return directions
}

func (d Direction) Count() int {
	return bits.OnesCount(uint(d))
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

func (d Direction) Add(dir Direction) Direction {
	return d | dir
}

func (d Direction) Remove(dir Direction) Direction {
	return d &^ dir
}

func (d Direction) Contains(dir Direction) bool {
	return d&dir != 0
}

func (d Direction) Symbol() string {
	if s, ok := symbols[d]; ok {
		return s
	}
	return " "
}
