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

// Emitters
//   - Emit on startup
//   - Emit every x
//   - Emit after x activation
//   - Random emit all port
//   - Emit on all port
//   - Emit on a random port
//   - Emit on one of the port round robin
//
// Triggers
// - Never (blocks signals)
// - Always
// - Fixed note
// - Random notes arpegiated (param range, algo)
type Node interface {
	Update(g *Grid, x, y int)
	Activated() bool
	Direction() Direction
	Symbol() string
	Color() string
	Reset()
}

type Trigger interface {
	Trig()
}

type Emitter interface {
	Emit()
}

type Signal struct {
	direction Direction
	updated   bool
}

func (s *Signal) Update(g *Grid, x, y int) {
	if s.updated {
		s.updated = false
	} else {
		g.Move(x, y, s.direction)
		s.updated = true
	}
}

func (s *Signal) Direction() Direction {
	return s.direction
}

func (s *Signal) Activated() bool {
	return true
}

func (s *Signal) Reset() {
	s.updated = false
}

func (s *Signal) Symbol() string {
	return " "
}

func (s *Signal) Color() string {
	return "15"
}

type BasicEmitter struct {
	direction Direction
	activated bool
	updated   bool
}

func (e *BasicEmitter) Emit() {
	e.activated = true
	e.updated = true
}

func (e *BasicEmitter) Update(g *Grid, x, y int) {
	if e.activated && !e.updated {
		g.Emit(x, y, e.direction)
		e.updated = true
		e.activated = false
	} else if e.updated {
		e.updated = false
	}
}

func (e *BasicEmitter) Direction() Direction {
	return e.direction
}

func (e *BasicEmitter) Activated() bool {
	return e.activated
}

func (e *BasicEmitter) Symbol() string {
	return "S"
}

func (e *BasicEmitter) Color() string {
	return "93"
}

func (e *BasicEmitter) Reset() {
	e.updated = false
}
