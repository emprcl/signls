package meta

import (
	"signls/core/common"
	"signls/midi"
)

type Command interface {
	// Active() bool
	// SetActive(active bool)
	// Copy() Command
	Execute()
	Executed() bool
}

type RootCommand struct {
	midi midi.Midi

	Value *common.ControlValue[uint8]

	executed bool
	active   bool
}

func (c *RootCommand) Executed() bool {
	return c.executed
}

func (c *RootCommand) Execute() {
	c.executed = true
}
