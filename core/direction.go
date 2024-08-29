package core

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

	// symbols maps Direction combinations to their corresponding string symbols.
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
// Returns a slice of Directions that make up the current Direction.
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

// Count returns the number of active directions in the current Direction.
// Utilizes the bits.OnesCount function to count the number of set bits.
func (d Direction) Count() int {
	return bits.OnesCount(uint(d))
}

// NextPosition calculates the new position (x, y) after moving in the specified Direction.
// The new coordinates are determined based on the direction of movement.
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

// Add combines the current Direction with another Direction using bitwise OR.
// This allows multiple directions to be represented in a single Direction.
func (d Direction) Add(dir Direction) Direction {
	return d | dir
}

// Remove subtracts a specific Direction from the current Direction using bitwise AND NOT.
// This effectively removes the specified direction from the current set of directions.
func (d Direction) Remove(dir Direction) Direction {
	return d &^ dir
}

// Contains checks if the current Direction includes a specific basic Direction.
// Returns true if the specified direction is part of the current Direction.
func (d Direction) Contains(dir Direction) bool {
	return d&dir != 0
}

// Symbol returns the string symbol associated with the current Direction.
// If a symbol is not defined for the current combination of directions, it returns a space.
func (d Direction) Symbol() string {
	if s, ok := symbols[d]; ok {
		return s
	}
	return " "
}
