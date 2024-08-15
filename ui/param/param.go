package param

import (
	"cykl/core"
)

type Values map[int]string

type Param interface {
	Name() string
	Value() int
	Display() string
	Set(value int)
	Increment()
	Decrement()
}

func NewParamsForNode(node core.Node) []Param {
	switch node.(type) {
	case *core.BangEmitter:
		return []Param{
			Direction{node: node},
			Key{node: node},
		}
	default:
		return []Param{}
	}
}
