package node

import (
	"cykl/core/common"
	"cykl/core/music"
)

type SignalContext struct {
	lastKey  music.Key
	position int
	arp      []int8
}

// Signal represents a directional pulse in a grid-based system.
// It contains the direction in which it is moving and the pulse value,
// which likely represents a timestamp or counter for synchronization.
type Signal struct {
	context   *SignalContext
	direction common.Direction // The current direction of the signal's movement.
	pulse     uint64           // The pulse value representing the last time the signal was updated.
}

// NewSignal creates a new Signal with the specified direction and pulse value.
// This function initializes the Signal with the provided parameters.
func NewSignal(direction common.Direction, pulse uint64) *Signal {
	return &Signal{
		context:   &SignalContext{},
		direction: direction,
		pulse:     pulse,
	}
}

// MustMove checks if the signal has to move in the grid.
func (s *Signal) MustMove(pulse uint64) bool {
	if !s.updated(pulse) {
		s.pulse = pulse
		return true
	}
	return false
}

// Direction returns the current direction of the signal's movement.
func (s *Signal) Direction() common.Direction {
	return s.direction
}

// SetDirection updates the direction of the signal's movement.
func (s *Signal) SetDirection(dir common.Direction) {
	s.direction = dir
}

// Activated is a placeholder method that currently always returns true.
// This could be extended to include more complex activation logic in the future.
func (s *Signal) Activated() bool {
	return true
}

// Symbol returns the symbol representing the signal.
func (s *Signal) Symbol() string {
	return "  "
}

// Name returns the name of the signal.
func (s *Signal) Name() string {
	return "signal"
}

// Color returns the color of the signal as a string.
func (s *Signal) Color() string {
	return "15"
}

// updated checks if the signal was updated on the given pulse.
func (s *Signal) updated(pulse uint64) bool {
	return s.pulse == pulse
}
