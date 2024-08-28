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
	Left()
	Right()
}

func NewParamsForNodes(grid *core.Grid, nodes []core.Node) []Param {
	if len(nodes) == 0 {
		return []Param{}
	}
	return []Param{
		Key{
			nodes: nodes,
			keys:  core.AllKeysInScale(grid.Key, grid.Scale),
			root:  grid.Key,
			mode:  KeyMode{nodes: nodes, modes: core.AllNoteBehaviors()},
		},
		Velocity{nodes: nodes},
		Length{nodes: nodes},
		Channel{nodes: nodes},
	}
}

func NewParamsForGrid(grid *core.Grid) []Param {
	return []Param{
		Root{grid: grid},
		Scale{grid: grid, scales: core.AllScales()},
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
