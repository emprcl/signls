package core

// Emitters
// - Never
// - Once on start
// - Every x pulses
// - Euclidean
//
// Triggers
// - Never (blocks signals)
// - Always
// - Fixed note
// - Random notes arpegiated (param range, algo)
type Node interface {
	Update(g *Grid, x, y int)
}

type Trigger interface {
	Trig()
}

type Signal struct {
	direction uint8
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

type OnceEmitter struct {
	direction  uint8
	shouldEmit bool
}

func (e *OnceEmitter) Update(g *Grid, x, y int) {
	if e.shouldEmit {
		g.Emit(x, y, e.direction)
	}
	e.shouldEmit = false
}

func (e *OnceEmitter) Direction() uint8 {
	return e.direction
}
