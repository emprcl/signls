package ui

import (
	"image/color"

	"signls/filesystem"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

// fallbackNodeColor is used when a node type has no color in the palette.
const fallbackNodeColor = "15"

// holeNodeKey is the node name whose color paints hole emitters and, since a
// teleport destination belongs to a hole, its destination marker too.
const holeNodeKey = "hole"

// styles holds every color-bearing lipgloss style the ui renders with, built
// once from a resolved theme palette. Layout-only styles that carry no color
// (controlStyle, cellStyle, ...) remain package-level vars.
type styles struct {
	grid    lipgloss.Style
	gridAlt lipgloss.Style
	// gridAltSet reports whether the alternating grid cells have their own
	// background; when false they fall through to the terminal background.
	gridAltSet      bool
	cursor          lipgloss.Style
	holeDestination lipgloss.Style
	selection       lipgloss.Style
	emitter         lipgloss.Style
	mutedEmitter    lipgloss.Style
	activeEmitter   lipgloss.Style
	activeCell      lipgloss.Style
	bank            lipgloss.Style
	bankOdd         lipgloss.Style
	activeBank      lipgloss.Style

	// inputCursor colors the text-input caret in the control bar.
	inputCursor color.Color
	// nodes maps a node's Name() to its themed color.
	nodes map[string]color.Color
}

// toColor turns a palette Color into a lipgloss color, using an adaptive color
// only when the light and dark variants differ.
func toColor(c filesystem.Color) color.Color {
	if c.Light == c.Dark || c.Dark == "" {
		return lipgloss.Color(c.Light)
	}
	return compat.AdaptiveColor{
		Light: lipgloss.Color(c.Light),
		Dark:  lipgloss.Color(c.Dark),
	}
}

// newStyles builds the render styles from a resolved palette.
func newStyles(p filesystem.Palette) styles {
	nodes := make(map[string]color.Color, len(p.Nodes))
	for k, v := range p.Nodes {
		nodes[k] = toColor(v)
	}
	// A teleport destination renders like the hole it belongs to: the hole
	// node color on a background, with the shared emitter foreground on top.
	holeColor, ok := nodes[holeNodeKey]
	if !ok {
		holeColor = lipgloss.Color(fallbackNodeColor)
	}
	return styles{
		grid: lipgloss.NewStyle().
			Background(toColor(p.GridBackground)),
		gridAlt: lipgloss.NewStyle().
			Background(toColor(p.GridBackgroundAlt)),
		gridAltSet: !p.GridBackgroundAlt.IsZero(),
		cursor: lipgloss.NewStyle().
			Background(toColor(p.Cursor)).
			Foreground(toColor(p.CursorForeground)),
		holeDestination: lipgloss.NewStyle().
			Background(holeColor).
			Foreground(toColor(p.EmitterForeground)),
		selection: lipgloss.NewStyle().
			Background(toColor(p.SelectionBackground)).
			Foreground(toColor(p.SelectionForeground)),
		emitter: lipgloss.NewStyle().
			Foreground(toColor(p.EmitterForeground)),
		mutedEmitter: lipgloss.NewStyle().
			Background(toColor(p.MutedEmitterBackground)).
			Foreground(toColor(p.MutedEmitterForeground)),
		activeEmitter: lipgloss.NewStyle().
			Background(toColor(p.ActiveEmitterBackground)).
			Foreground(toColor(p.ActiveEmitterForeground)),
		activeCell: cellStyle.
			Foreground(toColor(p.ActiveCell)),
		bank: lipgloss.NewStyle().
			MarginRight(1).
			Background(toColor(p.BankEven)).
			Foreground(toColor(p.BankForeground)),
		bankOdd: lipgloss.NewStyle().
			MarginRight(1).
			Background(toColor(p.BankOdd)).
			Foreground(toColor(p.BankForeground)),
		activeBank: lipgloss.NewStyle().
			MarginRight(1).
			Background(toColor(p.BankActive)).
			Foreground(toColor(p.BankForeground)),
		inputCursor: toColor(p.Cursor),
		nodes:       nodes,
	}
}

// nodeColor returns the themed color for a node type key, falling back to a
// sane default for unknown keys.
func (s styles) nodeColor(key string) color.Color {
	if c, ok := s.nodes[key]; ok {
		return c
	}
	return lipgloss.Color(fallbackNodeColor)
}
