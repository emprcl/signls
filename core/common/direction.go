package common

import "math/bits"

// Direction is a custom type representing directional values using bitwise operations.
type Direction int

// Constants representing individual directions as bit flags.
// Each direction is assigned a distinct bit position, allowing multiple directions
// to be combined using bitwise operations.
const (
	NONE Direction = 0
	UP   Direction = 1 << iota
	RIGHT
	DOWN
	LEFT
)

// Predefined variables for managing directions and their string representations.
var (
	// allDirections is a slice containing all the basic directional constants.
	allDirections = []Direction{UP, RIGHT, DOWN, LEFT}

	// symbols maps direction combinations to their corresponding string symbols.
	symbols = map[Direction]string{
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

// Decompose splits a Direction into its constituent basic directions.
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

// Count returns the number of active directions in the current direction.
func (d Direction) Count() int {
	return bits.OnesCount(uint(d))
}

// NextPosition calculates the new position (x, y) after moving in the specified direction
// from a given position.
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
		return 0, 0 // Default case, should not be reached if directions are properly handled.
	}
}

// Add combines the current direction with another direction.
func (d Direction) Add(dir Direction) Direction {
	return d | dir
}

// Remove subtracts a specific direction from the current direction.
func (d Direction) Remove(dir Direction) Direction {
	return d &^ dir
}

// Contains checks if the current direction includes a specific basic direction.
func (d Direction) Contains(dir Direction) bool {
	return d&dir != 0
}

// Symbol returns the string symbol associated with the current direction.
func (d Direction) Symbol() string {
	if s, ok := symbols[d]; ok {
		return s
	}
	return " "
}
