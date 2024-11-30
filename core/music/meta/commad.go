package meta

import (
	"signls/core/common"
)

type Command interface {
	Active() bool
	SetActive(active bool)
	Copy() Command
	Execute()
	Executed() bool
	Value() *common.ControlValue[int]
	Display() string
	Reset()
}
