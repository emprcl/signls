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
// Note behavior
// - Silence
// - Behavior from signal (how to?)
// - Fixed (tied to scale)
// - Random (tied to scale)
// - Random notes arpegiated (param range, algo) (tied to signal?)
type Node interface {
	Activated() bool
	Direction() Direction
	SetDirection(dir Direction)
	Symbol() string
	Name() string
	Color() string
}

type Movable interface {
	Move(g *Grid, x, y int)
}

type Emitter interface {
	Arm()
	Copy() Node
	Note() Note
	Muted() bool
	SetMute(mute bool)
	Trig(g *Grid, x, y int)
	Emit(g *Grid, x, y int)
	Reset()
}
