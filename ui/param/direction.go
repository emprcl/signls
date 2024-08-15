package param

import (
	"cykl/core"
	"fmt"
)

type Direction struct {
	node core.Node
}

func (d Direction) Name() string {
	return "direction"
}

func (d Direction) Display() string {
	return fmt.Sprintf("%s (%04b)", d.node.Direction().Symbol(), int(d.node.Direction()>>1))
}

func (d Direction) Value() int {
	return int(d.node.Direction())
}

func (d Direction) Increment() {
	// direction selection working differently
}

func (d Direction) Decrement() {
	// direction selection working differently
}

func (d Direction) Set(value int) {
	d.node.SetDirection(core.Direction(value))
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
	d.node.SetDirection(dir)
}