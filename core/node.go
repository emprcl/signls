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
	Activated() bool
	Direction() Direction
	SetDirection(dir Direction)
	Symbol() string
	Color() string
}

type Movable interface {
	Move(g *Grid, x, y int)
}

type Emitter interface {
	Arm()
	Note() Note
	Trig(g *Grid, x, y int)
	Emit(g *Grid, x, y int)
	Reset()
}
