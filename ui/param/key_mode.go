package param

import (
	"cykl/core"
)

type KeyMode struct {
	node  core.Node
	modes []core.NoteBehavior
}

func (k KeyMode) Name() string {
	return "mode"
}

func (k KeyMode) Display() string {
	return k.node.(*core.Emitter).Note().Behavior.Name()
}

func (k KeyMode) Value() int {
	return 0
}

func (k KeyMode) Increment() {
	k.Set(k.keyModeIndex() + 1)
}

func (k KeyMode) Decrement() {
	k.Set(k.keyModeIndex() - 1)
}

func (k KeyMode) Left() {}

func (k KeyMode) Right() {}

func (k KeyMode) Set(value int) {
	if value < 0 {
		value = len(k.modes) - 1
	} else if value >= len(k.modes) {
		value = 0
	}
	k.node.(*core.Emitter).Note().Behavior = k.modes[value]
}

func (k KeyMode) keyModeIndex() int {
	for i := 0; i < len(k.modes); i++ {
		if k.node.(*core.Emitter).Note().Behavior == k.modes[i] {
			return i
		}
	}
	return 0
}
