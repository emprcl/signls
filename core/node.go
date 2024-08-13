package core

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
	Armed() bool
	Direction() Direction
	Symbol() string
	Color() string
	Reset()
}

type Movable interface {
	Move(g *Grid, x, y int)
}

type Emitter interface {
	Arm()
	Note() Note
	Emit(g *Grid, x, y int)
	Trig(g *Grid, x, y int)
}
