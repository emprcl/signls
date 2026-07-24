package filesystem

import (
	"encoding/json"
	"sort"
)

// DefaultThemeName is the preset used when no (or an unknown) theme is selected.
const DefaultThemeName = "default"

// Color is a single theme color. It marshals to a plain JSON string ("190",
// "#ff8800") when its light and dark variants are identical, or to a
// {"light": ..., "dark": ...} object for colors that must adapt to the
// terminal background. Values are anything lipgloss understands: an ANSI index,
// an ANSI-256 index, or a hex string.
type Color struct {
	Light string
	Dark  string
}

// solid builds a Color that renders the same on light and dark backgrounds.
func solid(v string) Color { return Color{Light: v, Dark: v} }

// adaptive builds a Color that adapts to the terminal background.
func adaptive(light, dark string) Color { return Color{Light: light, Dark: dark} }

// IsZero reports whether the color is unset (used to detect overrides).
func (c Color) IsZero() bool { return c.Light == "" && c.Dark == "" }

// MarshalJSON writes a bare string for solid colors and an object for adaptive
// ones, keeping user-facing config files as terse as possible.
func (c Color) MarshalJSON() ([]byte, error) {
	if c.Light == c.Dark {
		return json.Marshal(c.Light)
	}
	return json.Marshal(struct {
		Light string `json:"light"`
		Dark  string `json:"dark"`
	}{c.Light, c.Dark})
}

// UnmarshalJSON accepts either a bare string (solid color) or a
// {"light","dark"} object (adaptive color).
func (c *Color) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		c.Light, c.Dark = s, s
		return nil
	}
	var obj struct {
		Light string `json:"light"`
		Dark  string `json:"dark"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	c.Light, c.Dark = obj.Light, obj.Dark
	return nil
}

// Palette holds every color the ui renders with. Nodes maps a node's Name()
// (e.g. "bang", "dice") to the color used for that node type. The "hole" entry
// also colors the teleport-destination marker, so both share a single color.
type Palette struct {
	Cursor Color `json:"cursor"`
	// GridBackground colors one set of the alternating empty grid cells, and
	// GridBackgroundAlt the other. Leave GridBackgroundAlt unset to let those
	// cells fall through to the terminal background (the original look).
	GridBackground          Color            `json:"gridBackground"`
	GridBackgroundAlt       Color            `json:"gridBackgroundAlt"`
	SelectionBackground     Color            `json:"selectionBackground"`
	SelectionForeground     Color            `json:"selectionForeground"`
	EmitterForeground       Color            `json:"emitterForeground"`
	ActiveEmitterBackground Color            `json:"activeEmitterBackground"`
	ActiveEmitterForeground Color            `json:"activeEmitterForeground"`
	MutedEmitterBackground  Color            `json:"mutedEmitterBackground"`
	MutedEmitterForeground  Color            `json:"mutedEmitterForeground"`
	ActiveCell              Color            `json:"activeCell"`
	BankForeground          Color            `json:"bankForeground"`
	BankEven                Color            `json:"bankEven"`
	BankOdd                 Color            `json:"bankOdd"`
	BankActive              Color            `json:"bankActive"`
	Nodes                   map[string]Color `json:"nodes"`
}

// MarshalJSON writes only the colors that are actually set. Palettes are stored
// as overrides layered over a preset, so an unset field must stay absent from
// the file rather than being written as an empty string (which json's
// omitempty cannot do for struct-valued fields).
func (p Palette) MarshalJSON() ([]byte, error) {
	opt := func(c Color) *Color {
		if c.IsZero() {
			return nil
		}
		return &c
	}
	var nodes map[string]Color
	for k, v := range p.Nodes {
		if v.IsZero() {
			continue
		}
		if nodes == nil {
			nodes = make(map[string]Color, len(p.Nodes))
		}
		nodes[k] = v
	}
	return json.Marshal(struct {
		Cursor                  *Color           `json:"cursor,omitempty"`
		GridBackground          *Color           `json:"gridBackground,omitempty"`
		GridBackgroundAlt       *Color           `json:"gridBackgroundAlt,omitempty"`
		SelectionBackground     *Color           `json:"selectionBackground,omitempty"`
		SelectionForeground     *Color           `json:"selectionForeground,omitempty"`
		EmitterForeground       *Color           `json:"emitterForeground,omitempty"`
		ActiveEmitterBackground *Color           `json:"activeEmitterBackground,omitempty"`
		ActiveEmitterForeground *Color           `json:"activeEmitterForeground,omitempty"`
		MutedEmitterBackground  *Color           `json:"mutedEmitterBackground,omitempty"`
		MutedEmitterForeground  *Color           `json:"mutedEmitterForeground,omitempty"`
		ActiveCell              *Color           `json:"activeCell,omitempty"`
		BankForeground          *Color           `json:"bankForeground,omitempty"`
		BankEven                *Color           `json:"bankEven,omitempty"`
		BankOdd                 *Color           `json:"bankOdd,omitempty"`
		BankActive              *Color           `json:"bankActive,omitempty"`
		Nodes                   map[string]Color `json:"nodes,omitempty"`
	}{
		Cursor:                  opt(p.Cursor),
		GridBackground:          opt(p.GridBackground),
		GridBackgroundAlt:       opt(p.GridBackgroundAlt),
		SelectionBackground:     opt(p.SelectionBackground),
		SelectionForeground:     opt(p.SelectionForeground),
		EmitterForeground:       opt(p.EmitterForeground),
		ActiveEmitterBackground: opt(p.ActiveEmitterBackground),
		ActiveEmitterForeground: opt(p.ActiveEmitterForeground),
		MutedEmitterBackground:  opt(p.MutedEmitterBackground),
		MutedEmitterForeground:  opt(p.MutedEmitterForeground),
		ActiveCell:              opt(p.ActiveCell),
		BankForeground:          opt(p.BankForeground),
		BankEven:                opt(p.BankEven),
		BankOdd:                 opt(p.BankOdd),
		BankActive:              opt(p.BankActive),
		Nodes:                   nodes,
	})
}

// merge returns a copy of p with every non-zero field of o applied on top. It
// is how user overrides are layered over a preset palette.
func (p Palette) merge(o Palette) Palette {
	// Clone the node map so we never mutate the preset stored in presets().
	nodes := make(map[string]Color, len(p.Nodes))
	for k, v := range p.Nodes {
		nodes[k] = v
	}
	p.Nodes = nodes

	set := func(dst *Color, src Color) {
		if !src.IsZero() {
			*dst = src
		}
	}
	set(&p.Cursor, o.Cursor)
	set(&p.GridBackground, o.GridBackground)
	set(&p.GridBackgroundAlt, o.GridBackgroundAlt)
	set(&p.SelectionBackground, o.SelectionBackground)
	set(&p.SelectionForeground, o.SelectionForeground)
	set(&p.EmitterForeground, o.EmitterForeground)
	set(&p.ActiveEmitterBackground, o.ActiveEmitterBackground)
	set(&p.ActiveEmitterForeground, o.ActiveEmitterForeground)
	set(&p.MutedEmitterBackground, o.MutedEmitterBackground)
	set(&p.MutedEmitterForeground, o.MutedEmitterForeground)
	set(&p.ActiveCell, o.ActiveCell)
	set(&p.BankForeground, o.BankForeground)
	set(&p.BankEven, o.BankEven)
	set(&p.BankOdd, o.BankOdd)
	set(&p.BankActive, o.BankActive)
	for k, v := range o.Nodes {
		if !v.IsZero() {
			p.Nodes[k] = v
		}
	}
	return p
}

// ThemeConfig is the serialized theme selection: the name of a built-in preset
// plus an optional palette whose set fields override the preset.
type ThemeConfig struct {
	Name    string   `json:"name"`
	Palette *Palette `json:"palette,omitempty"`
}

// MarshalJSON drops an all-empty palette override so the config file never
// grows an empty "palette": {} block; a nil or fully-unset palette is omitted.
func (t ThemeConfig) MarshalJSON() ([]byte, error) {
	type alias ThemeConfig // sheds the method set, avoiding recursion
	a := alias(t)
	if a.Palette != nil {
		if b, err := json.Marshal(a.Palette); err == nil && string(b) == "{}" {
			a.Palette = nil
		}
	}
	return json.Marshal(a)
}

// Resolve returns the full palette: the selected preset (default when the name
// is empty or unknown) with any user overrides merged on top.
func (t ThemeConfig) Resolve() Palette {
	base, ok := presets()[t.Name]
	if !ok {
		base = presets()[DefaultThemeName]
	}
	if t.Palette != nil {
		base = base.merge(*t.Palette)
	}
	return base
}

// HasTheme reports whether name refers to a built-in preset.
func HasTheme(name string) bool {
	_, ok := presets()[name]
	return ok
}

// ThemeNames returns the sorted names of every built-in preset.
func ThemeNames() []string {
	names := make([]string, 0, len(presets()))
	for name := range presets() {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// presets returns the built-in palettes. It rebuilds them on each call so
// callers can never mutate shared state.
func presets() map[string]Palette {
	return map[string]Palette{
		DefaultThemeName: defaultPalette(),
		"mono":           monoPalette(),
		"neon":           neonPalette(),
	}
}

// defaultPalette is the original signls color scheme.
func defaultPalette() Palette {
	return Palette{
		Cursor:                  solid("190"),
		GridBackground:          adaptive("254", "234"),
		SelectionBackground:     solid("238"),
		SelectionForeground:     solid("244"),
		EmitterForeground:       solid("15"),
		ActiveEmitterBackground: adaptive("0", "15"),
		ActiveEmitterForeground: adaptive("15", "0"),
		MutedEmitterBackground:  solid("247"),
		MutedEmitterForeground:  solid("236"),
		ActiveCell:              solid("190"),
		BankForeground:          solid("0"),
		BankEven:                solid("79"),
		BankOdd:                 solid("85"),
		BankActive:              solid("15"),
		Nodes: map[string]Color{
			"signal": solid("15"),
			"bang":   solid("165"),
			"dice":   solid("33"),
			"cycle":  solid("63"),
			"euclid": solid("162"),
			"pass":   solid("35"),
			"hole":   solid("124"),
			"spread": solid("56"),
			"zone":   solid("197"),
			"toll":   solid("39"),
		},
	}
}

// monoPalette is a restrained grayscale theme with a single amber accent.
//
// Node colors are used both as cell backgrounds (with a bright symbol on top)
// and, when a node fires, as the symbol color on a light highlight. They are
// therefore kept in the dark-gray band (234-244): bright symbols read clearly
// on them, and they read clearly on the light active-cell highlight. Both
// terminal backgrounds are covered because every emitter cell paints its own
// background rather than relying on the terminal's.
func monoPalette() Palette {
	return Palette{
		Cursor:                  solid("223"),
		GridBackground:          adaptive("253", "235"),
		GridBackgroundAlt:       adaptive("255", "233"),
		SelectionBackground:     solid("240"),
		SelectionForeground:     solid("253"),
		EmitterForeground:       solid("231"),
		ActiveEmitterBackground: solid("252"),
		ActiveEmitterForeground: solid("235"),
		MutedEmitterBackground:  solid("240"),
		MutedEmitterForeground:  solid("250"),
		ActiveCell:              solid("223"),
		BankForeground:          solid("233"),
		BankEven:                solid("245"),
		BankOdd:                 solid("250"),
		BankActive:              solid("223"),
		Nodes: map[string]Color{
			"signal": solid("231"),
			"bang":   solid("244"),
			"dice":   solid("242"),
			"cycle":  solid("240"),
			"euclid": solid("238"),
			"pass":   solid("236"),
			"hole":   solid("243"),
			"spread": solid("241"),
			"zone":   solid("239"),
			"toll":   solid("237"),
		},
	}
}

// neonPalette is a vivid, high-contrast theme tuned for dark terminals.
func neonPalette() Palette {
	return Palette{
		Cursor:                  solid("201"),
		GridBackground:          adaptive("237", "234"),
		GridBackgroundAlt:       adaptive("236", "232"),
		SelectionBackground:     solid("54"),
		SelectionForeground:     solid("213"),
		EmitterForeground:       solid("231"),
		ActiveEmitterBackground: adaptive("57", "51"),
		ActiveEmitterForeground: adaptive("231", "16"),
		MutedEmitterBackground:  solid("238"),
		MutedEmitterForeground:  solid("244"),
		ActiveCell:              solid("201"),
		BankForeground:          solid("16"),
		BankEven:                solid("51"),
		BankOdd:                 solid("45"),
		BankActive:              solid("231"),
		Nodes: map[string]Color{
			"signal": solid("231"),
			"bang":   solid("201"),
			"dice":   solid("51"),
			"cycle":  solid("141"),
			"euclid": solid("213"),
			"pass":   solid("48"),
			"hole":   solid("196"),
			"spread": solid("99"),
			"zone":   solid("205"),
			"toll":   solid("39"),
		},
	}
}
