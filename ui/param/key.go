package param

import (
	"time"

	"cykl/core/common"
	"cykl/core/music"
)

type Key struct {
	nodes []common.Node
	keys  []music.Key
	root  music.Key
}

func (k Key) Name() string {
	return k.nodes[0].(music.Audible).Note().KeyName()
}

func (k Key) Display() string {
	return k.nodes[0].(music.Audible).Note().KeyName()
}

func (k Key) Value() int {
	return int(k.nodes[0].(music.Audible).Note().KeyValue())
}

func (k Key) Increment() {
	k.Set(k.keyIndex() + 1)
}

func (k Key) Decrement() {
	k.Set(k.keyIndex() + -1)
}

func (k Key) Left() {
	//k.mode.Left()
}

func (k Key) Right() {
	//k.mode.Right()
}

func (k Key) Set(value int) {
	for _, n := range k.nodes {
		n.(music.Audible).Note().SetKey(k.keys[value], k.root)
	}
}

func (k Key) Preview() {
	go func() {
		n := *k.nodes[0].(music.Audible).Note()
		n.Play(music.Key(60), music.CHROMATIC)
		time.Sleep(300 * time.Millisecond)
		n.Stop()
	}()
}

func (k Key) keyIndex() int {
	for i := 0; i < len(k.keys); i++ {
		if k.nodes[0].(music.Audible).Note().KeyValue() == k.keys[i] {
			return i
		}
	}
	return 0
}
