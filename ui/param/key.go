package param

import (
	"cykl/core"
)

type Key struct {
	node core.Node
}

func (k Key) Name() string {
	return "key"
}

func (k Key) Display() string {
	return k.node.(core.Emitter).Note().Name()
}

func (k Key) Value() int {
	return int(k.node.(core.Emitter).Note().Key)
}

func (k Key) Increment() {
	k.Set(k.Value() + 1)
}

func (k Key) Decrement() {
	k.Set(k.Value() - 1)
}

func (k Key) Set(value int) {
	k.node.(core.Emitter).Note().SetKey(uint8(value))
}
