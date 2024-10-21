package node

import (
	"signls/core/common"
)

const (
	HoleDestinationSymbol = "H+"
)

type HoleEmitter struct {
	activated    int
	originX      int
	originY      int
	DestinationX *common.ControlValue[int]
	DestinationY *common.ControlValue[int]
}

func NewHoleEmitter(direction common.Direction, x, y, width, height int) *HoleEmitter {
	return &HoleEmitter{
		originX:      x,
		originY:      y,
		DestinationX: common.NewControlValue[int](x, 0, width-1),
		DestinationY: common.NewControlValue[int](y, 0, height-1),
	}
}

func (e *HoleEmitter) Copy(dx, dy int) common.Node {
	newDestinationX := *e.DestinationX
	newDestinationY := *e.DestinationY
	newDestinationX.Set(newDestinationX.Value() - e.originX + dx)
	newDestinationY.Set(newDestinationY.Value() - e.originY + dy)
	return &HoleEmitter{
		originX:      dx,
		originY:      dy,
		DestinationX: &newDestinationX,
		DestinationY: &newDestinationY,
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
	return e.DestinationX.Computed(), e.DestinationY.Computed()
}

func (e *HoleEmitter) Destination() (int, int) {
	return e.DestinationX.Value(), e.DestinationY.Value()
}

func (e *HoleEmitter) DestinationAmount() (int, int) {
	return e.DestinationX.RandomAmount(), e.DestinationY.RandomAmount()
}

func (e *HoleEmitter) SetDestination(x, y int) {
	e.DestinationX.Set(x)
	e.DestinationY.Set(y)
}

func (e *HoleEmitter) SetDestinationAmount(x, y int) {
	e.DestinationX.SetRandomAmount(x)
	e.DestinationY.SetRandomAmount(y)
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
