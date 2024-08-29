package core

// Signal represents a directional pulse in a grid-based system.
// It contains the direction in which it is moving and the pulse value,
// which likely represents a timestamp or counter for synchronization.
type Signal struct {
	direction Direction // The current direction of the signal's movement.
	pulse     uint64    // The pulse value representing the last time the signal was updated.
}

// NewSignal creates a new Signal with the specified direction and pulse value.
// This function initializes the Signal with the provided parameters.
func NewSignal(direction Direction, pulse uint64) *Signal {
	return &Signal{
		direction: direction,
		pulse:     pulse,
	}
}

// Move attempts to move the signal in its current direction on the grid.
// The movement only occurs if the signal has not already been updated during the current pulse.
func (s *Signal) Move(g *Grid, x, y int) {
	if !s.updated(g.pulse) {
		g.Move(x, y, s.direction)
		s.pulse = g.pulse
	}
}

// Direction returns the current direction of the signal's movement.
func (s *Signal) Direction() Direction {
	return s.direction
}

// SetDirection updates the direction of the signal's movement.
func (s *Signal) SetDirection(dir Direction) {
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

// updated checks if the signal has already been updated during the current pulse.
// If the signal's pulse matches the provided pulse, it returns true.
func (s Signal) updated(pulse uint64) bool {
	return s.pulse == pulse
}
