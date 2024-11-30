package param

import (
	"signls/core/common"
	"signls/core/field"
	"signls/core/music"
	"signls/core/node"
	"signls/core/theory"
)

const (
	defaultControlParamsNumber = 8
)

type Param interface {
	Name() string
	Value() int
	AltValue() int
	Display() string
	Set(value int)
	SetAlt(value int)
	SetEditValue(input string)

	Up()
	Down()
	Left()
	Right()
	AltUp()
	AltDown()
	AltLeft()
	AltRight()
}

func NewParamsForNodes(grid *field.Grid, nodes []common.Node) [][]Param {
	if len(nodes) == 0 {
		return [][]Param{}
	}

	if isHomogeneousNode[*node.HoleEmitter](nodes) {
		return [][]Param{
			{
				Destination{
					nodes:  nodes,
					width:  grid.Width,
					height: grid.Height,
				},
			},
		}
	} else if isHomogeneousBehavior[*node.TollEmitter](nodes) {
		return [][]Param{
			append(
				DefaultEmitterParams(grid, nodes),
				Threshold{nodes: nodes},
			),
			DefaultEmitterControlChanges(nodes),
			DefaultEmitterMetaCommands(nodes),
		}
	} else if isHomogeneousNode[*node.EuclidEmitter](nodes) {
		return [][]Param{
			append(
				DefaultEmitterParams(grid, nodes),
				Steps{nodes: nodes},
				Triggers{nodes: nodes},
				Offset{nodes: nodes},
			),
			DefaultEmitterControlChanges(nodes),
			DefaultEmitterMetaCommands(nodes),
		}
	}

	emitters := filterNodes[music.Audible](nodes)

	return [][]Param{
		DefaultEmitterParams(grid, emitters),
		DefaultEmitterControlChanges(emitters),
		DefaultEmitterMetaCommands(emitters),
	}
}

func DefaultEmitterParams(grid *field.Grid, nodes []common.Node) []Param {
	keyMode := KeyModeRandom
	if nodes[0].(music.Audible).Note().Key.IsSilent() {
		keyMode = KeyModeSilent
	}
	return []Param{
		&Key{
			nodes: nodes,
			keys:  theory.AllKeysInScale(grid.Key, grid.Scale),
			root:  grid.Key,
			scale: grid.Scale,
			mode:  keyMode,
		},
		Velocity{nodes: nodes},
		Length{nodes: nodes},
		Channel{nodes: nodes},
		Probability{nodes: nodes},
	}
}

func DefaultEmitterControlChanges(nodes []common.Node) []Param {
	params := make([]Param, defaultControlParamsNumber)
	for i := range params {
		params[i] = CC{index: i, nodes: nodes}
	}
	return params
}

func DefaultEmitterMetaCommands(nodes []common.Node) []Param {
	return []Param{
		RootCmd{nodes: nodes},
		ScaleCmd{nodes: nodes},
	}
}

func NewParamsForGrid(grid *field.Grid) []Param {
	return []Param{
		Root{grid: grid},
		Scale{grid: grid, scales: theory.AllScales()},
	}
}

func NewParamsForMidi(grid *field.Grid) [][]Param {
	return [][]Param{
		{
			ClockSend{grid: grid},
			TransportSend{grid: grid},
			DefaultDevice{grid: grid},
		},
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
