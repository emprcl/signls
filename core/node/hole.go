package node

import (
	"cykl/core/common"
)

type HoleEmitter struct {
	activated    int
	destinationX int
	destinationY int
}

func NewHoleEmitter(direction common.Direction, x, y int) *HoleEmitter {
	return &HoleEmitter{
		destinationX: x,
		destinationY: y,
	}
}

func (e *HoleEmitter) Copy(dx, dy int) common.Node {
	return &HoleEmitter{
		destinationX: e.destinationX + dx,
		destinationY: e.destinationY + dy,
	}
}

func (e *HoleEmitter) Activated() bool {
	return e.activated > 0
}

func (s *HoleEmitter) Direction() common.Direction {
	return common.NONE
}

func (s *HoleEmitter) SetDirection(dir common.Direction) {}

func (e *HoleEmitter) Teleport() (int, int) {
	e.activated = common.PulsesPerStep + 1
	return e.Destination()
}

func (e *HoleEmitter) Destination() (int, int) {
	return e.destinationX, e.destinationY
}

func (e *HoleEmitter) SetDestination(x, y int) {
	e.destinationX, e.destinationY = x, y
}

func (e *HoleEmitter) Tick() {
	if e.activated <= 0 {
		return
	}
	e.activated--
}

func (e *HoleEmitter) Reset() {
	e.activated = 0
}

func (s *HoleEmitter) Symbol() string {
	return "H⬢"
}

func (s *HoleEmitter) Name() string {
	return "hole"
}

func (s *HoleEmitter) Color() string {
	return "124"
}