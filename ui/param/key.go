package param

import (
	"cykl/core"
	"time"
)

type Key struct {
	node core.Node
}

func (k Key) Name() string {
	return "key"
}

func (k Key) Display() string {
	return k.node.(core.Emitter).Note().KeyName()
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

func (k Key) Preview() {
	go func() {
		n := *k.node.(core.Emitter).Note()
		n.Play()
		time.Sleep(time.Second)
		n.Stop()
	}()
}
