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
	case *core.BangEmitter, *core.SpreadEmitter:
		return []Param{
			Direction{node: node},
			Key{node: node},
			Velocity{node: node},
			Length{node: node},
			Channel{node: node},
		}
	default:
		return []Param{}
	}
}

func NewParamsForGrid(grid *core.Grid) []Param {
	return []Param{
		Root{grid: grid},
	}
}

func Get(name string, params []Param) Param {
	for _, p := range params {
		if p.Name() == name {
			return p
		}
	}
	return params[0]
}
