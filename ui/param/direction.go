package param

import (
	"signls/core/common"
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

func (d Direction) AltValue() int {
	return 0
}

func (d Direction) Up() {}

func (d Direction) Down() {}

func (d Direction) Left() {}

func (d Direction) Right() {}

func (d Direction) AltUp() {}

func (d Direction) AltDown() {}

func (d Direction) AltLeft() {}

func (d Direction) AltRight() {}

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
