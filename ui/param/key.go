package param

import (
	"fmt"
	"time"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/core/node"
)

type Key struct {
	nodes []common.Node
	keys  []music.Key
	root  music.Key
	mode  KeyMode
}

func (k Key) Name() string {
	return "key"
}

func (k Key) Display() string {
	switch k.mode.Display() {
	case "silent":
		return "•"
	case "random":
		return fmt.Sprintf("%s%s", "r", k.nodes[0].(*node.BaseEmitter).Note().KeyName())
	default:
		return k.nodes[0].(*node.BaseEmitter).Note().KeyName()
	}
}

func (k Key) Value() int {
	return int(k.nodes[0].(*node.BaseEmitter).Note().KeyValue())
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
	for _, n := range k.nodes {
		n.(*node.BaseEmitter).Note().SetKey(k.keys[value], k.root)
	}
}

func (k Key) Preview() {
	go func() {
		n := *k.nodes[0].(*node.BaseEmitter).Note()
		n.Play(music.Key(60), music.CHROMATIC)
		time.Sleep(300 * time.Millisecond)
		n.Stop()
	}()
}

func (k Key) keyIndex() int {
	for i := 0; i < len(k.keys); i++ {
		if k.nodes[0].(*node.BaseEmitter).Note().KeyValue() == k.keys[i] {
			return i
		}
	}
	return 0
}
