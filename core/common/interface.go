package common

// Node represents a general interface for any element that can exist on the grid.
// It defines common behavior for all types of nodes, such as emitters or signals.
type Node interface {
	// Activated returns whether the node is currently active.
	Activated() bool

	// Direction returns the current direction the node is facing.
	Direction() Direction

	// SetDirection sets the direction for the node.
	SetDirection(dir Direction)

	// Symbol returns a string that visually represents the node.
	Symbol() string

	// Name returns the name of the node, typically used for identifying its type.
	Name() string

	// Color returns a string representing the color code for the node,
	// which could be used for rendering in a terminal or GUI.
	Color() string
}

// EmitterBehavior defines the behavior of different types of emitters.
type EmitterBehavior interface {
	// ArmedOnStart indicates if the emitter is armed when the grid starts.
	ArmedOnStart() bool

	// Copy makes a copy of the behavior.
	Copy() EmitterBehavior

	// EmitDirections determines which directions the emitter will emit signals
	// based on its current direction, the incoming direction, and the pulse count.
	EmitDirections(dir Direction, inDir Direction, pulse uint64) Direction

	// ShouldPropagate indicates if triggers should be propagated to direct
	// neighbors.
	ShouldPropagate() bool

	// Reset resets the behavior state.
	Reset()

	// Symbol returns a string representation of the emitter
	Symbol() string

	// Name returns the name of the emitter type.
	Name() string

	// Color returns the color code associated with the emitter.
	Color() string
}

// Movable represents an interface for nodes that can move within the grid.
// This could be implemented by any node that changes position over time.
type Movable interface {
	// MustMove checks if the node should move for the current pulse.
	MustMove(pulse uint64) bool
}

type Behavioral interface {
	Behavior() EmitterBehavior
	SetBehavior(behavior EmitterBehavior)
}

type Tickable interface {
	Tick()
	Reset()
}

type Copyable interface {
	Copy(dx, dy int) Node
}

type Parameter[T any] interface {
	Value() T
	Computed() T
	Last() T
	Set(value T)
	RandomAmount() int
	SetRandomAmount(amount int)
}
