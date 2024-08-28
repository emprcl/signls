package param

import (
	"cykl/core"
)

type Direction struct {
	nodes []core.Node
}

func NewDirection(nodes []core.Node) Direction {
	return Direction{
		nodes: nodes,
	}
}

func (d Direction) Name() string {
	return "direction"
}

func (d Direction) Display() string {
	return d.nodes[0].Direction().Symbol()
}

func (d Direction) Value() int {
	return int(d.nodes[0].Direction())
}

func (d Direction) Increment() {
	// direction selection working differently
}

func (d Direction) Decrement() {
	// direction selection working differently
}

func (d Direction) Left() {}

func (d Direction) Right() {}

func (d Direction) Set(value int) {
	for _, node := range d.nodes {
		node.SetDirection(core.Direction(value))
	}
}

func (d Direction) SetFromKeyString(key string) {
	var dir core.Direction
	switch key {
	case "up":
		dir = core.UP
	case "right":
		dir = core.RIGHT
	case "down":
		dir = core.DOWN
	case "left":
		dir = core.LEFT
	default:
		dir = core.UP
	}
	for _, node := range d.nodes {
		node.SetDirection(dir)
	}
}
