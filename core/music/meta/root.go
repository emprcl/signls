package meta

import (
	"signls/core/common"
	"signls/midi"
)

const (
	defaultKey = 60 // Middle C
	maxKey     = 127
	minKey     = 21
)

type RootCommand struct {
	value    *common.ControlValue[int]
	executed bool
	active   bool
}

func NewRootCommand() *RootCommand {
	return &RootCommand{
		value: common.NewControlValue[int](defaultKey, minKey, maxKey),
	}
}

func (c *RootCommand) Copy() Command {
	newValue := *c.value
	return &RootCommand{
		value:  &newValue,
		active: c.active,
	}
}

func (c *RootCommand) Active() bool {
	return c.active
}

func (c *RootCommand) SetActive(active bool) {
	c.active = active
}

func (c *RootCommand) Executed() bool {
	return c.executed
}

func (c *RootCommand) Execute() {
	if !c.active {
		return
	}
	c.executed = true
}

func (c *RootCommand) Value() *common.ControlValue[int] {
	return c.value
}

func (c *RootCommand) Display() string {
	return midi.Note(uint8(c.value.Value()))
}

func (c *RootCommand) Name() string {
	return "root"
}

func (c *RootCommand) Reset() {
	c.executed = false
}
