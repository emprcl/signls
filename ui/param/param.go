package param

import (
	"cykl/core/common"
	"cykl/core/field"
	"cykl/core/music"
	"cykl/core/node"
)

type Param interface {
	Name() string
	Value() int
	AltValue() int
	Display() string
	Set(value int)
	SetAlt(value int)

	Up()
	Down()
	Left()
	Right()
	AltUp()
	AltDown()
	AltLeft()
	AltRight()
}

func NewParamsForNodes(grid *field.Grid, nodes []common.Node) []Param {
	if len(nodes) == 0 {
		return []Param{}
	}

	if isHomogeneousNode[*node.HoleEmitter](nodes) {
		return []Param{
			Destination{
				nodes:  nodes,
				width:  grid.Width,
				height: grid.Height,
			},
		}
	} else if isHomogeneousBehavior[*node.TollEmitter](nodes) {
		return append(
			DefaultEmitterParams(grid, nodes),
			Threshold{nodes: nodes},
		)
	} else if isHomogeneousNode[*node.EuclidEmitter](nodes) {
		return append(
			DefaultEmitterParams(grid, nodes),
			Steps{nodes: nodes},
			Triggers{nodes: nodes},
			Offset{nodes: nodes},
		)
	}

	emitters := filterNodes[music.Audible](nodes)
	return DefaultEmitterParams(grid, emitters)
}

func DefaultEmitterParams(grid *field.Grid, nodes []common.Node) []Param {
	return []Param{
		&Key{
			nodes: nodes,
			keys:  music.AllKeysInScale(grid.Key, grid.Scale),
			root:  grid.Key,
		},
		Velocity{nodes: nodes},
		Length{nodes: nodes},
		Channel{nodes: nodes},
		Probability{nodes: nodes},
	}
}

func NewParamsForGrid(grid *field.Grid) []Param {
	return []Param{
		Root{grid: grid},
		Scale{grid: grid, scales: music.AllScales()},
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

func filterNodes[T any](nodes []common.Node) []common.Node {
	filteredNodes := []common.Node{}
	for _, n := range nodes {
		if _, ok := n.(T); !ok {
			continue
		}
		filteredNodes = append(filteredNodes, n)
	}
	return filteredNodes
}

func isHomogeneousNode[T any](nodes []common.Node) bool {
	for _, n := range nodes {
		if _, ok := n.(T); !ok {
			return false
		}
	}
	return true
}

func isHomogeneousBehavior[T any](nodes []common.Node) bool {
	for _, n := range nodes {
		if _, ok := n.(common.Behavioral); !ok {
			return false
		}

		if _, ok := n.(common.Behavioral).Behavior().(T); !ok {
			return false
		}
	}
	return true
}
