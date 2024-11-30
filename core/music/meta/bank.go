package meta

import (
	"fmt"
	"signls/core/common"
)

const (
	defaultGrid = 0
	maxGrid     = 31
	minGrid     = 0
)

type BankCommand struct {
	value    *common.ControlValue[int]
	executed bool
	active   bool
}

func NewBankCommand() *BankCommand {
	return &BankCommand{
		value: common.NewControlValue[int](defaultGrid, minGrid, maxGrid),
	}
}

func (c *BankCommand) Copy() Command {
	newValue := *c.value
	return &BankCommand{
		value:  &newValue,
		active: c.active,
	}
}

func (c *BankCommand) Active() bool {
	return c.active
}

func (c *BankCommand) SetActive(active bool) {
	c.active = active
}

func (c *BankCommand) Executed() bool {
	return c.executed
}

func (c *BankCommand) Execute() {
	if !c.active {
		return
	}
	c.executed = true
}

func (c *BankCommand) Value() *common.ControlValue[int] {
	return c.value
}

func (c *BankCommand) Display() string {
	return fmt.Sprintf("%d", c.value.Value()+1)
}

func (c *BankCommand) Name() string {
	return "bank"
}

func (c *BankCommand) Reset() {
	c.executed = false
}
