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

type BasicEmitter struct {
	Direction  uint8
	shouldEmit bool
}

func (e *BasicEmitter) Emit() {
	e.shouldEmit = true
}

func (e *BasicEmitter) Update(g *Grid, x, y int) {
	if e.shouldEmit {
		g.Emit(x, y, e.Direction)
	}
	e.shouldEmit = false
}
