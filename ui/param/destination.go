package param

import (
	"fmt"
	"strconv"
	"strings"

	"signls/core/common"
	"signls/core/node"
)

type Destination struct {
	nodes  []common.Node
	width  int
	height int
}

func (d Destination) Name() string {
	return "dest"
}

func (s Destination) Help() string {
	return ""
}

func (d Destination) Display() string {
	x, y := d.nodes[0].(*node.HoleEmitter).Destination()
	amountX, amountY := d.nodes[0].(*node.HoleEmitter).DestinationAmount()
	if amountX == 0 && amountY == 0 {
		return fmt.Sprintf("%d,%d", x, y)
	}
	return fmt.Sprintf("%d%+d,%d%+d", x, amountX, y, amountY)
}

func (d Destination) Value() int {
	return 0
}

func (d Destination) AltValue() int {
	return 0
}

func (d Destination) Position() (int, int) {
	return d.nodes[0].(*node.HoleEmitter).Destination()
}

func (d Destination) Up() {
	d.SetDestination(0, -1)
}

func (d Destination) Down() {
	d.SetDestination(0, 1)
}

func (d Destination) Left() {
	d.SetDestination(-1, 0)
}

func (d Destination) Right() {
	d.SetDestination(1, 0)
}

func (d Destination) AltUp() {
	d.SetDestinationAmount(0, 1)
}

func (d Destination) AltDown() {
	d.SetDestinationAmount(0, -1)
}

func (d Destination) AltLeft() {
	d.SetDestinationAmount(-1, 0)
}

func (d Destination) AltRight() {
	d.SetDestinationAmount(1, 0)
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

func (d Destination) SetDestinationAmount(dAmountX, dAmountY int) {
	for _, n := range d.nodes {
		amountX, amountY := n.(*node.HoleEmitter).DestinationAmount()
		n.(*node.HoleEmitter).SetDestinationAmount(amountX+dAmountX, amountY+dAmountY)
	}
}

func (d Destination) SetEditValue(input string) {
	coordinates := strings.Split(input, ",")
	if len(coordinates) != 2 {
		return
	}
	x, errX := strconv.Atoi(coordinates[0])
	y, errY := strconv.Atoi(coordinates[1])
	if errX != nil || errY != nil {
		return
	}
	for _, n := range d.nodes {
		n.(*node.HoleEmitter).SetDestination(x, y)
	}
}
