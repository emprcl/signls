package param

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/node"
)

type Destination struct {
	nodes  []common.Node
	width  int
	height int
}

func (d Destination) Name() string {
	return "dest"
}

func (d Destination) Display() string {
	x, y := d.nodes[0].(*node.HoleEmitter).Destination()
	return fmt.Sprintf("%d,%d", x, y)
}

func (d Destination) Value() int {
	return 0
}

func (d Destination) Position() (int, int) {
	return d.nodes[0].(*node.HoleEmitter).Destination()
}

func (d Destination) Increment() {
	d.SetDestination(0, -1)
}

func (d Destination) Decrement() {
	d.SetDestination(0, 1)
}

func (d Destination) Left() {
	d.SetDestination(-1, 0)
}

func (d Destination) Right() {
	d.SetDestination(1, 0)
}

func (d Destination) Set(value int) {}

func (d Destination) SetAlt(value int) {}

func (d Destination) SetDestination(dx, dy int) {
	for _, n := range d.nodes {
		x, y := n.(*node.HoleEmitter).Destination()
		destinationX := x + dx
		destinationY := y + dy
		if destinationX < 0 ||
			destinationX >= d.width ||
			destinationY < 0 ||
			destinationY >= d.height {
			continue
		}
		n.(*node.HoleEmitter).SetDestination(x+dx, y+dy)
	}
}
