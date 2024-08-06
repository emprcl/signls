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
	Reset()
}

type Trigger interface {
	Trig()
}

type Emitter interface {
	Emit()
}

type Signal struct {
	Direction uint8
	updated   bool
}

func (s *Signal) Update(g *Grid, x, y int) {
	if s.updated {
		s.updated = false
	} else {
		g.Move(x, y, s.Direction)
		s.updated = true
	}
}

func (s *Signal) Reset() {
	s.updated = false
}

type BasicEmitter struct {
	Direction uint8
	Activated bool
	updated   bool
}

func (e *BasicEmitter) Emit() {
	e.Activated = true
	e.updated = true
}

func (e *BasicEmitter) Update(g *Grid, x, y int) {
	if e.Activated && !e.updated {
		g.Emit(x, y, e.Direction)
		e.updated = true
		e.Activated = false
	} else if e.updated {
		e.updated = false
	}
}

func (e *BasicEmitter) Reset() {
	e.updated = false
}
