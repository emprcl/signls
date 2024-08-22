package param

import (
	"cykl/core"
	"fmt"
	"time"
)

type Key struct {
	node core.Node
	keys []core.Key
	root core.Key
	mode KeyMode
}

func (k Key) Name() string {
	return "key"
}

func (k Key) Display() string {
	switch k.mode.Display() {
	case "silent":
		return "•"
	case "random":
		return fmt.Sprintf("%s%s", "r", k.node.(*core.Emitter).Note().KeyName())
	default:
		return k.node.(*core.Emitter).Note().KeyName()
	}
}

func (k Key) Value() int {
	return int(k.node.(*core.Emitter).Note().KeyValue())
}

func (k Key) Increment() {
	k.Set(k.keyIndex() + 1)
}

func (k Key) Decrement() {
	k.Set(k.keyIndex() - 1)
}

func (k Key) Left() {
	k.mode.Decrement()
}

func (k Key) Right() {
	k.mode.Increment()
}

func (k Key) Set(value int) {
	k.node.(*core.Emitter).Note().SetKey(k.keys[value], k.root)
}

func (k Key) Preview() {
	go func() {
		n := *k.node.(*core.Emitter).Note()
		n.Play(core.Key(60), core.CHROMATIC)
		time.Sleep(300 * time.Millisecond)
		n.Stop()
	}()
}

func (k Key) keyIndex() int {
	for i := 0; i < len(k.keys); i++ {
		if k.node.(*core.Emitter).Note().KeyValue() == k.keys[i] {
			return i
		}
	}
	return 0
}
