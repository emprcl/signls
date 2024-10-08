package node

import (
	"cykl/core/common"
)

const (
	HoleDestinationSymbol = "H+"
)

type HoleEmitter struct {
	activated    int
	originX      int
	originY      int
	destinationX *common.ControlValue[int]
	destinationY *common.ControlValue[int]
}

func NewHoleEmitter(direction common.Direction, x, y, width, height int) *HoleEmitter {
	return &HoleEmitter{
		originX:      x,
		originY:      y,
		destinationX: common.NewControlValue[int](x, 0, width),
		destinationY: common.NewControlValue[int](y, 0, height),
	}
}

func (e *HoleEmitter) Copy(dx, dy int) common.Node {
	newDestinationX := *e.destinationX
	newDestinationY := *e.destinationY
	newDestinationX.Set(newDestinationX.Value() - e.originX + dx)
	newDestinationY.Set(newDestinationY.Value() - e.originY + dy)
	return &HoleEmitter{
		originX:      dx,
		originY:      dy,
		destinationX: &newDestinationX,
		destinationY: &newDestinationY,
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
	return e.destinationX.Computed(), e.destinationY.Computed()
}

func (e *HoleEmitter) Destination() (int, int) {
	return e.destinationX.Value(), e.destinationY.Value()
}

func (e *HoleEmitter) DestinationAmount() (int, int) {
	return e.destinationX.RandomAmount(), e.destinationY.RandomAmount()
}

func (e *HoleEmitter) SetDestination(x, y int) {
	e.destinationX.Set(x)
	e.destinationY.Set(y)
}

func (e *HoleEmitter) SetDestinationAmount(x, y int) {
	e.destinationX.SetRandomAmount(x)
	e.destinationY.SetRandomAmount(y)
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
