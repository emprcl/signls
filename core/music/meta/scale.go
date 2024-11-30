package meta

import (
	"signls/core/common"
	"signls/core/theory"
)

const (
	defaultScale = 0
	minScale     = 0
)

type ScaleCommand struct {
	value    *common.ControlValue[int]
	executed bool
	active   bool
}

func NewScaleCommand() *ScaleCommand {
	return &ScaleCommand{
		value: common.NewControlValue[int](defaultScale, minScale, len(theory.AllScales())-1),
	}
}

func (c *ScaleCommand) Copy() Command {
	newValue := *c.value
	return &ScaleCommand{
		value:  &newValue,
		active: c.active,
	}
}

func (c *ScaleCommand) Active() bool {
	return c.active
}

func (c *ScaleCommand) SetActive(active bool) {
	c.active = active
}

func (c *ScaleCommand) Executed() bool {
	return c.executed
}

func (c *ScaleCommand) Execute() {
	if !c.active {
		return
	}
	c.executed = true
}

func (c *ScaleCommand) Value() *common.ControlValue[int] {
	return c.value
}

func (c *ScaleCommand) Display() string {
	return theory.AllScales()[c.Value().Value()].Name()
}

func (c *ScaleCommand) Reset() {
	c.executed = false
}
