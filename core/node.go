package core

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

// Movable represents an interface for nodes that can move within the grid.
// This could be implemented by any node that changes position over time.
type Movable interface {
	// Move is a method that takes a reference to the grid, and the current x, y
	// position of the node. It defines how the node should move on the grid.
	Move(g *Grid, x, y int)
}
