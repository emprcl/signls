package meta

import (
	"fmt"
	"signls/core/common"
)

const (
	defaultTempo = 120
	maxTempo     = 300
	minTempo     = 1
)

type TempoCommand struct {
	value    *common.ControlValue[int]
	executed bool
	active   bool
}

func NewTempoCommand() *TempoCommand {
	return &TempoCommand{
		value: common.NewControlValue[int](defaultTempo, minTempo, maxTempo),
	}
}

func (c *TempoCommand) Copy() Command {
	newValue := *c.value
	return &TempoCommand{
		value:  &newValue,
		active: c.active,
	}
}

func (c *TempoCommand) Active() bool {
	return c.active
}

func (c *TempoCommand) SetActive(active bool) {
	c.active = active
}

func (c *TempoCommand) Executed() bool {
	return c.executed
}

func (c *TempoCommand) Execute() {
	if !c.active {
		return
	}
	c.executed = true
}

func (c *TempoCommand) Value() *common.ControlValue[int] {
	return c.value
}

func (c *TempoCommand) Display() string {
	return fmt.Sprintf("%d", c.value.Value())
}

func (c *TempoCommand) Name() string {
	return "tempo"
}

func (c *TempoCommand) Reset() {
	c.executed = false
}
