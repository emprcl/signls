package param

import (
	"cykl/core/common"
)

type Direction struct {
	nodes []common.Node
}

func NewDirection(nodes []common.Node) Direction {
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
		node.SetDirection(common.Direction(value))
	}
}

func (d Direction) SetAlt(value int) {}

func (d Direction) SetFromKeyString(key string) {
	var dir common.Direction
	switch key {
	case "up":
		dir = common.UP
	case "right":
		dir = common.RIGHT
	case "down":
		dir = common.DOWN
	case "left":
		dir = common.LEFT
	default:
		dir = common.UP
	}
	for _, node := range d.nodes {
		node.SetDirection(dir)
	}
}
