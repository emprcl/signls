package filesystem

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// DefaultThemeName is the preset used when no theme is selected.
const DefaultThemeName = "default"

// Color is a single theme color. It marshals to a plain JSON string ("190",
// "#ff8800") when its light and dark variants are identical, or to a
// {"light": ..., "dark": ...} object for colors that must adapt to the
// terminal background. Values are an ANSI index (0-255) or a hex string.
type Color struct {
	Light string
	Dark  string
}

// solid builds a Color that renders the same on light and dark backgrounds.
func solid(v string) *Color { return &Color{Light: v, Dark: v} }

// adaptive builds a Color that adapts to the terminal background.
func adaptive(light, dark string) *Color { return &Color{Light: light, Dark: dark} }

// IsZero reports whether the color carries no value at all.
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
// {"light","dark"} object (adaptive color). Anything else is reported rather
// than half-decoded, so that `"cursor": 190` (a number instead of a string)
// names its own mistake.
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
		return fmt.Errorf(
			"invalid color %s: expected a string (\"190\", \"#ff0088\") or an object ({\"light\": ..., \"dark\": ...})",
			data,
		)
	}
	c.Light, c.Dark = obj.Light, obj.Dark
	return nil
}

// Palette holds every color the ui renders with. A nil field means "this
// palette does not set that color": a preset may leave one out, and a user
// palette sets only the colors it overrides.
//
// Adding a color is a one-line change. Marshalling (omitempty), merging,
// validation and the set of keys accepted from a config file all derive from
// this struct by walking its Color fields, so there is no parallel list that
// can silently fall out of sync.
type Palette struct {
	Cursor           *Color `json:"cursor,omitempty"`
	CursorForeground *Color `json:"cursorForeground,omitempty"`
	// GridBackground colors one set of the alternating empty grid cells, and
	// GridBackgroundAlt the other. Leave GridBackgroundAlt unset to let those
	// cells fall through to the terminal background (the original look).
	GridBackground          *Color `json:"gridBackground,omitempty"`
	GridBackgroundAlt       *Color `json:"gridBackgroundAlt,omitempty"`
	SelectionBackground     *Color `json:"selectionBackground,omitempty"`
	SelectionForeground     *Color `json:"selectionForeground,omitempty"`
	EmitterForeground       *Color `json:"emitterForeground,omitempty"`
	ActiveEmitterBackground *Color `json:"activeEmitterBackground,omitempty"`
	ActiveEmitterForeground *Color `json:"activeEmitterForeground,omitempty"`
	MutedEmitterBackground  *Color `json:"mutedEmitterBackground,omitempty"`
	MutedEmitterForeground  *Color `json:"mutedEmitterForeground,omitempty"`
	ActiveCell              *Color `json:"activeCell,omitempty"`
	BankForeground          *Color `json:"bankForeground,omitempty"`
	BankEven                *Color `json:"bankEven,omitempty"`
	BankOdd                 *Color `json:"bankOdd,omitempty"`
	BankActive              *Color `json:"bankActive,omitempty"`
	// BankActiveForeground labels the active grid's slot. Leave it unset to
	// label it like every other slot; a theme needs it only when BankActive is
	// too light or too dark for BankForeground to read on.
	BankActiveForeground *Color `json:"bankActiveForeground,omitempty"`
	// Nodes maps a node's Name() (e.g. "bang", "dice") to the color used for
	// that node type. The "hole" entry also colors the teleport-destination
	// marker, so both share a single color.
	Nodes map[string]*Color `json:"nodes,omitempty"`
}

// nodeKeys are the node types a palette colors: every node's Name(). A preset
// is expected to cover all of them, since a node with no color of its own is a
// node the theme cannot draw.
var nodeKeys = []string{
	"signal", "bang", "dice", "cycle", "euclid",
	"pass", "hole", "spread", "zone", "toll",
}

// sameNodes gives every node type one color, for presets that tell node types
// apart by their glyph rather than by their color.
func sameNodes(v string) map[string]*Color {
	nodes := make(map[string]*Color, len(nodeKeys))
	for _, name := range nodeKeys {
		nodes[name] = solid(v)
	}
	return nodes
}

// colorPtrType is the type of every themeable color field of Palette.
var colorPtrType = reflect.TypeOf((*Color)(nil))

// colorField is one themeable color: the json key it is written under, and a
// pointer to the field holding it.
type colorField struct {
	key string
	ptr **Color
}

// colorFields returns every color field of the palette in declaration order.
// It is the single source of truth behind merge and Validate.
func (p *Palette) colorFields() []colorField {
	v := reflect.ValueOf(p).Elem()
	t := v.Type()
	fields := make([]colorField, 0, t.NumField())
	for i := range t.NumField() {
		if v.Field(i).Type() != colorPtrType {
			continue // the Nodes map, handled separately
		}
		fields = append(fields, colorField{
			key: jsonKey(t.Field(i)),
			ptr: v.Field(i).Addr().Interface().(**Color),
		})
	}
	return fields
}

// jsonKey returns the json object key a struct field is written under.
func jsonKey(f reflect.StructField) string {
	key, _, _ := strings.Cut(f.Tag.Get("json"), ",")
	return key
}

// paletteKeys returns every key a palette accepts in a config file, sorted.
func paletteKeys() []string {
	t := reflect.TypeOf(Palette{})
	keys := make([]string, 0, t.NumField())
	for i := range t.NumField() {
		keys = append(keys, jsonKey(t.Field(i)))
	}
	sort.Strings(keys)
	return keys
}

// UnmarshalJSON rejects keys the palette does not know. Such a key would
// otherwise be dropped on load and then erased from the file on the next save,
// so a mistyped key would silently do nothing and then disappear.
func (p *Palette) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	known := paletteKeys()
	unknown := make([]string, 0, len(raw))
	for key := range raw {
		if !slices.Contains(known, key) {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return fmt.Errorf(
			"unknown palette key %q (valid keys: %s)",
			unknown[0], strings.Join(known, ", "),
		)
	}
	type alias Palette // sheds the method set, avoiding recursion
	return json.Unmarshal(data, (*alias)(p))
}

// Validate reports the first problem found in a palette. Colors are checked at
// load time because an unusable value is not visibly an error: lipgloss drops
// what it cannot parse, so a mistyped hex string or an out-of-range ANSI index
// renders as an invisible cursor or a missing node rather than as a complaint.
func (p Palette) Validate() error {
	for _, f := range p.colorFields() {
		if err := validateColor(f.key, *f.ptr); err != nil {
			return err
		}
	}
	names := make([]string, 0, len(p.Nodes))
	for name := range p.Nodes {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !slices.Contains(nodeKeys, name) {
			return fmt.Errorf(
				"unknown node %q in palette (valid nodes: %s)",
				name, strings.Join(NodeNames(), ", "),
			)
		}
		if err := validateColor("nodes."+name, p.Nodes[name]); err != nil {
			return err
		}
	}
	return nil
}

// validateColor checks one palette entry, naming the key it belongs to.
func validateColor(key string, c *Color) error {
	if c == nil {
		return nil
	}
	if c.Light == "" || c.Dark == "" {
		return fmt.Errorf(
			"palette color %q is empty: remove the key to keep the theme's own color",
			key,
		)
	}
	for _, v := range []string{c.Light, c.Dark} {
		if err := validateColorValue(v); err != nil {
			return fmt.Errorf("palette color %q: %w", key, err)
		}
	}
	return nil
}

// validateColorValue accepts what lipgloss can actually render: an ANSI index
// (0-255) or a hex string (#rgb or #rrggbb).
func validateColorValue(v string) error {
	if hex, ok := strings.CutPrefix(v, "#"); ok {
		if len(hex) != 3 && len(hex) != 6 {
			return fmt.Errorf("invalid hex color %q: expected #rgb or #rrggbb", v)
		}
		if strings.TrimLeft(strings.ToLower(hex), "0123456789abcdef") != "" {
			return fmt.Errorf("invalid hex color %q: not a hex number", v)
		}
		return nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 || n > 255 {
		return fmt.Errorf(
			"invalid color %q: expected an ANSI index 0-255 or a hex string like \"#ff0088\"",
			v,
		)
	}
	return nil
}

// merge returns a copy of p with every color o sets applied on top. It is how
// user overrides are layered over a preset palette.
func (p Palette) merge(o Palette) Palette {
	// Clone the node map so we never mutate the preset stored in presets().
	nodes := make(map[string]*Color, len(p.Nodes))
	for k, v := range p.Nodes {
		nodes[k] = v
	}
	p.Nodes = nodes

	dst := p.colorFields()
	for i, f := range o.colorFields() {
		if *f.ptr == nil {
			continue
		}
		c := **f.ptr // copy, so the merged palette shares nothing with o
		*dst[i].ptr = &c
	}
	for k, v := range o.Nodes {
		if v == nil {
			continue
		}
		c := *v
		p.Nodes[k] = &c
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

// Validate reports the first problem found in the theme selection. An empty
// name is valid: it means the config predates themes and gets the default.
func (t ThemeConfig) Validate() error {
	if t.Name != "" && !HasTheme(t.Name) {
		return fmt.Errorf(
			"unknown theme %q (available themes: %s)",
			t.Name, strings.Join(ThemeNames(), ", "),
		)
	}
	if t.Palette == nil {
		return nil
	}
	return t.Palette.Validate()
}

// Resolve returns the full palette: the selected preset with any user overrides
// merged on top. Callers are expected to have validated the config first; an
// unknown name resolves to the default preset rather than to no colors at all.
func (t ThemeConfig) Resolve() Palette {
	base, ok := presets()[t.Name]
	if !ok {
		base = presets()[DefaultThemeName]
	}
	if t.Palette != nil {
		base = base.merge(*t.Palette)
	}
	// Colors that default to another color rather than to nothing are filled in
	// here, so the ui reads one complete palette instead of re-deriving them.
	if base.BankActiveForeground == nil {
		base.BankActiveForeground = base.BankForeground
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

// NodeNames returns the sorted node names a palette can color.
func NodeNames() []string {
	names := slices.Clone(nodeKeys)
	sort.Strings(names)
	return names
}

// presets returns the built-in palettes. It rebuilds them on each call so
// callers can never mutate shared state.
func presets() map[string]Palette {
	return map[string]Palette{
		DefaultThemeName: defaultPalette(),
		"mono":           monoPalette(),
		"pink":           pinkPalette(),
	}
}

// defaultPalette is the original signls color scheme.
func defaultPalette() Palette {
	return Palette{
		Cursor:                  solid("190"),
		CursorForeground:        solid("0"),
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
		Nodes: map[string]*Color{
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

// monoPalette is a three-color theme: black, white, and a grey that carries the
// states white alone cannot keep apart. The grid is a black canvas and idle
// nodes are white glyphs on it. Brightness then ranks state by urgency: the
// cursor is a white block (the loudest thing on screen, and the only white
// one), a firing node is a grey block, and a muted node recedes to a grey glyph
// on the canvas: still legible, plainly quieter than its unmuted neighbours.
// Node types share one color, so they are told apart by their glyph.
func monoPalette() Palette {
	const (
		black = "#000000"
		white = "#ffffff"
		grey  = "#808080"
	)
	return Palette{
		Cursor:            solid(white),
		CursorForeground:  solid(black),
		GridBackground:    solid(black),
		GridBackgroundAlt: solid(black),
		// Selected empty cells are a grey block, distinct from the cursor's
		// white one.
		SelectionBackground: solid(grey),
		SelectionForeground: solid(white),
		EmitterForeground:   solid(white),
		// A firing node keeps the inversion of the default theme, in grey so it
		// can never be mistaken for the cursor.
		ActiveEmitterBackground: solid(grey),
		ActiveEmitterForeground: solid(white),
		// Muted nodes sit directly on the canvas: no block, dimmed glyph.
		MutedEmitterBackground: solid(black),
		MutedEmitterForeground: solid(grey),
		ActiveCell:             solid(white),
		// Bank cells are the canvas color so the white-block selector (the
		// cursor) stands out against them; the active grid is a grey block.
		BankForeground: solid(white),
		BankEven:       solid(black),
		BankOdd:        solid(black),
		BankActive:     solid(grey),
		Nodes:          sameNodes(black),
	}
}

// pinkPalette mirrors mono's structure with a deep pink canvas: white glyphs
// for idle nodes, a white block for the cursor, and a pale pink standing in for
// mono's grey: a pale block for a firing node, a pale glyph for a muted one.
// Text over the pale pink is drawn in the canvas pink rather than white, since
// white on pale pink does not read.
func pinkPalette() Palette {
	const (
		pink  = "#c2185b"
		pale  = "#f8bbd0"
		white = "#ffffff"
	)
	return Palette{
		Cursor:                  solid(white),
		CursorForeground:        solid(pink),
		GridBackground:          solid(pink),
		GridBackgroundAlt:       solid(pink),
		SelectionBackground:     solid(pale),
		SelectionForeground:     solid(pink),
		EmitterForeground:       solid(white),
		ActiveEmitterBackground: solid(pale),
		ActiveEmitterForeground: solid(pink),
		MutedEmitterBackground:  solid(pink),
		MutedEmitterForeground:  solid(pale),
		ActiveCell:              solid(white),
		// Bank cells are the canvas color so the white-block selector (the
		// cursor) stands out against them; the active grid is a pale block,
		// labelled in pink since white does not read on pale pink.
		BankForeground:       solid(white),
		BankEven:             solid(pink),
		BankOdd:              solid(pink),
		BankActive:           solid(pale),
		BankActiveForeground: solid(pink),
		Nodes:                sameNodes(pink),
	}
}
