package filesystem

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestPresetsAreComplete guards the reflective walk against a color or a node
// added to Palette but forgotten in a preset. A missing color renders unstyled
// and a missing node renders in the ui's fallback color, neither of which any
// preset chose.
func TestPresetsAreComplete(t *testing.T) {
	for _, name := range ThemeNames() {
		palette := ThemeConfig{Name: name}.Resolve()
		for _, f := range palette.colorFields() {
			// The default theme deliberately leaves the alternating grid cells
			// on the terminal background.
			if f.key == "gridBackgroundAlt" && name == DefaultThemeName {
				continue
			}
			if *f.ptr == nil {
				t.Errorf("theme %q: color %q is not set", name, f.key)
			}
		}
		for _, node := range nodeKeys {
			if palette.Nodes[node] == nil {
				t.Errorf("theme %q: node %q has no color", name, node)
			}
		}
		if err := palette.Validate(); err != nil {
			t.Errorf("theme %q: %v", name, err)
		}
	}
}

// TestPresetsKeepStatesApart is the reason mono and pink carry a third color.
// Each of these pairs is a distinction the player reads off the grid at a
// glance (where the cursor is, what is firing, what is muted), and a theme that
// paints two of them with the same colors hides one of them completely.
//
// The states are grouped by what they draw, because only states drawing the
// same thing can be confused for one another: a firing node and a selected
// empty cell may share their colors, since one carries the node's glyph and the
// other the selection's dots.
func TestPresetsKeepStatesApart(t *testing.T) {
	// state names a rendered cell by the two palette colors it draws with.
	type state struct {
		name   string
		bg, fg *Color
	}
	for _, theme := range ThemeNames() {
		p := ThemeConfig{Name: theme}.Resolve()
		groups := map[string][]state{
			"cells holding a node": {
				{name: "cursor", bg: p.Cursor, fg: p.CursorForeground},
				{name: "firing node", bg: p.ActiveEmitterBackground, fg: p.ActiveEmitterForeground},
				{name: "muted node", bg: p.MutedEmitterBackground, fg: p.MutedEmitterForeground},
				{name: "idle node", bg: p.Nodes["bang"], fg: p.EmitterForeground},
			},
			"empty cells": {
				{name: "cursor", bg: p.Cursor, fg: p.CursorForeground},
				{name: "selection", bg: p.SelectionBackground, fg: p.SelectionForeground},
				{name: "grid", bg: p.GridBackground, fg: p.GridBackground},
			},
			"bank slots": {
				{name: "cursor", bg: p.Cursor, fg: p.CursorForeground},
				{name: "active grid", bg: p.BankActive, fg: p.BankActiveForeground},
				{name: "other grids", bg: p.BankEven, fg: p.BankForeground},
			},
		}
		for group, states := range groups {
			for i, a := range states {
				for _, b := range states[i+1:] {
					if a.bg == nil || b.bg == nil {
						continue
					}
					if *a.bg == *b.bg && *a.fg == *b.fg {
						t.Errorf(
							"theme %q, %s: %s and %s render identically (%s on %s)",
							theme, group, a.name, b.name, a.fg.Light, a.bg.Light,
						)
					}
				}
			}
		}
	}
}

func TestThemeConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr string
	}{
		{
			name:   "preset only",
			config: `{"name": "mono"}`,
		},
		{
			name:   "no name predates themes",
			config: `{}`,
		},
		{
			name:   "valid overrides",
			config: `{"name": "default", "palette": {"cursor": "205", "gridBackground": {"light": "255", "dark": "232"}, "nodes": {"bang": "#ff0088"}}}`,
		},
		{
			name:    "unknown preset",
			config:  `{"name": "nord"}`,
			wantErr: `unknown theme "nord"`,
		},
		{
			name:    "mistyped palette key",
			config:  `{"name": "default", "palette": {"cursorForground": "205"}}`,
			wantErr: `unknown palette key "cursorForground"`,
		},
		{
			name:    "mistyped node",
			config:  `{"name": "default", "palette": {"nodes": {"dise": "33"}}}`,
			wantErr: `unknown node "dise"`,
		},
		{
			name:    "ansi color out of range",
			config:  `{"name": "default", "palette": {"cursor": "300"}}`,
			wantErr: `invalid color "300"`,
		},
		{
			name:    "truncated hex color",
			config:  `{"name": "default", "palette": {"cursor": "#12345"}}`,
			wantErr: `invalid hex color "#12345"`,
		},
		{
			name:    "hex color with a non-hex digit",
			config:  `{"name": "default", "palette": {"cursor": "#gg0088"}}`,
			wantErr: `invalid hex color "#gg0088"`,
		},
		{
			name:    "empty color",
			config:  `{"name": "default", "palette": {"cursor": ""}}`,
			wantErr: `palette color "cursor" is empty`,
		},
		{
			name:    "half-written adaptive color",
			config:  `{"name": "default", "palette": {"cursor": {"light": "205"}}}`,
			wantErr: `palette color "cursor" is empty`,
		},
		{
			name:    "number instead of a color string",
			config:  `{"name": "default", "palette": {"cursor": 205}}`,
			wantErr: `invalid color 205`,
		},
		{
			name:    "node color out of range",
			config:  `{"name": "default", "palette": {"nodes": {"bang": "999"}}}`,
			wantErr: `palette color "nodes.bang"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config ThemeConfig
			err := json.Unmarshal([]byte(tt.config), &config)
			if err == nil {
				err = config.Validate()
			}
			switch {
			case tt.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tt.wantErr == "":
			case err == nil:
				t.Fatalf("no error, want one containing %q", tt.wantErr)
			case !strings.Contains(err.Error(), tt.wantErr):
				t.Fatalf("error = %q, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

// TestMergeOverridesPresetWithoutMutatingIt covers what a user palette is for:
// the colors it sets win, the ones it omits keep the preset's value, and the
// preset itself is left alone for the next caller.
func TestMergeOverridesPresetWithoutMutatingIt(t *testing.T) {
	config := ThemeConfig{
		Name: "mono",
		Palette: &Palette{
			Cursor: solid("205"),
			Nodes:  map[string]*Color{"bang": solid("9")},
		},
	}

	merged := config.Resolve()
	if got := merged.Cursor.Light; got != "205" {
		t.Errorf("cursor = %q, want the override %q", got, "205")
	}
	if got, want := merged.CursorForeground.Light, monoPalette().CursorForeground.Light; got != want {
		t.Errorf("cursorForeground = %q, want the preset's %q", got, want)
	}
	if got := merged.Nodes["bang"].Light; got != "9" {
		t.Errorf("nodes.bang = %q, want the override %q", got, "9")
	}
	if got, want := merged.Nodes["dice"].Light, monoPalette().Nodes["dice"].Light; got != want {
		t.Errorf("nodes.dice = %q, want the preset's %q", got, want)
	}

	preset := monoPalette()
	if got := preset.Cursor.Light; got == "205" {
		t.Error("merge leaked the override into the mono preset")
	}
	if got := preset.Nodes["bang"].Light; got == "9" {
		t.Error("merge leaked the node override into the mono preset")
	}
}

// TestPaletteRoundTrip checks that a palette written back to the config file
// keeps only what was set: unset colors must stay absent rather than reappear
// as empty strings that later fail validation.
func TestPaletteRoundTrip(t *testing.T) {
	config := ThemeConfig{Name: "pink", Palette: &Palette{Cursor: solid("205")}}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"name":"pink","palette":{"cursor":"205"}}`; string(data) != want {
		t.Fatalf("Marshal = %s, want %s", data, want)
	}

	var back ThemeConfig
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if err := back.Validate(); err != nil {
		t.Fatalf("Validate after round trip: %v", err)
	}
}

// TestEmptyPaletteIsNotPersisted keeps an empty "palette": {} block out of the
// config file a first run writes.
func TestEmptyPaletteIsNotPersisted(t *testing.T) {
	data, err := json.Marshal(ThemeConfig{Name: DefaultThemeName, Palette: &Palette{}})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"name":"default"}`; string(data) != want {
		t.Errorf("Marshal = %s, want %s", data, want)
	}
}
