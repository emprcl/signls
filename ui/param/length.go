package param

import (
	"cykl/core"
	"fmt"
)

type Length struct {
	node core.Node
}

func (l Length) Name() string {
	return "length"
}

func (l Length) Display() string {
	return fmt.Sprintf("%d", l.node.(core.Emitter).Note().Length)
}

func (l Length) Value() int {
	return int(l.node.(core.Emitter).Note().Length)
}

func (l Length) Increment() {
	l.Set(l.Value() + 1)
}

func (l Length) Decrement() {
	l.Set(l.Value() - 1)
}

func (l Length) Set(value int) {
	l.node.(core.Emitter).Note().SetLength(uint8(value))
}
