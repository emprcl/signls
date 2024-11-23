package filesystem

// KeyMap represents a keyboard mapping loaded from a json file.
type KeyMap struct {
	Play string `json:"play"`

	Up    string `json:"up"`
	Right string `json:"right"`
	Down  string `json:"down"`
	Left  string `json:"left"`

	SelectionUp    string `json:"selection_up"`
	SelectionRight string `json:"selection_right"`
	SelectionDown  string `json:"selection_down"`
	SelectionLeft  string `json:"selection_left"`

	EditUp    string `json:"edit_up"`
	EditRight string `json:"edit_right"`
	EditDown  string `json:"edit_down"`
	EditLeft  string `json:"edit_left"`

	EditInput string `json:"edit_input"`

	Bank string `json:"bank"`

	AddBang   string `json:"add_bang"`
	AddEuclid string `json:"add_euclid"`
	AddPass   string `json:"add_pass"`
	AddSpread string `json:"add_spread"`
	AddCycle  string `json:"add_cycle"`
	AddDice   string `json:"add_dice"`
	AddToll   string `json:"add_toll"`
	AddZone   string `json:"add_zone"`
	AddHole   string `json:"add_hole"`

	Copy  string `json:"copy"`
	Cut   string `json:"cut"`
	Paste string `json:"paste"`

	EditNode    string `json:"edit_node"`
	RemoveNode  string `json:"remove_node"`
	TriggerNode string `json:"trigger_node"`

	MuteNode    string `json:"mute_node"`
	MuteAllNode string `json:"mute_all_node"`

	RootNoteUp   string `json:"root_note_up"`
	RootNoteDown string `json:"root_note_down"`
	ScaleUp      string `json:"scale_up"`
	ScaleDown    string `json:"scale_down"`
	TempoUp      string `json:"tempo_up"`
	TempoDown    string `json:"tempo_down"`

	Configuration   string `json:"configuration"`
	FitGridToWindow string `json:"fit_grid_to_window"`

	Cancel string `json:"cancel"`

	Help string `json:"help"`
	Quit string `json:"quit"`
}

// NewDefaultAzertyKeyMap returns a new default KeyMap for azerty keyboards.
func NewDefaultAzertyKeyMap() KeyMap {
	return KeyMap{
		Play: " ",

		Up:    "up",
		Right: "right",
		Down:  "down",
		Left:  "left",

		SelectionUp:    "shift+up",
		SelectionRight: "shift+right",
		SelectionDown:  "shift+down",
		SelectionLeft:  "shift+left",

		EditUp:    "ctrl+up",
		EditRight: "ctrl+right",
		EditDown:  "ctrl+down",
		EditLeft:  "ctrl+left",

		EditInput: ":",

		Bank: "tab",

		AddBang:   "&",
		AddEuclid: "é",
		AddPass:   "\"",
		AddSpread: "'",
		AddCycle:  "(",
		AddDice:   "-",
		AddToll:   "è",
		AddZone:   "_",
		AddHole:   "ç",

		Copy:  "ctrl+c",
		Cut:   "ctrl+x",
		Paste: "ctrl+v",

		EditNode:    "enter",
		RemoveNode:  "backspace",
		TriggerNode: "!",

		MuteNode:    "m",
		MuteAllNode: "M",

		RootNoteUp:   "*",
		RootNoteDown: "ù",
		ScaleUp:      "µ",
		ScaleDown:    "%",
		TempoUp:      "=",
		TempoDown:    ")",

		Configuration:   "f2",
		FitGridToWindow: "f10",

		Cancel: "esc",

		Help: "?",
		Quit: "ctrl+q",
	}
}

// NewDefaultAzertyMacKeyMap returns a new default KeyMap for azerty mac
// keyboards.
func NewDefaultAzertyMacKeyMap() KeyMap {
	return KeyMap{
		Play: " ",

		Up:    "up",
		Right: "right",
		Down:  "down",
		Left:  "left",

		SelectionUp:    "shift+up",
		SelectionRight: "shift+right",
		SelectionDown:  "shift+down",
		SelectionLeft:  "shift+left",

		EditUp:    "ctrl+up",
		EditRight: "ctrl+right",
		EditDown:  "ctrl+down",
		EditLeft:  "ctrl+left",

		EditInput: ":",

		Bank: "tab",

		AddBang:   "&",
		AddEuclid: "é",
		AddPass:   "\"",
		AddSpread: "'",
		AddCycle:  "(",
		AddDice:   "§",
		AddToll:   "è",
		AddZone:   "!",
		AddHole:   "ç",

		Copy:  "ctrl+c",
		Cut:   "ctrl+x",
		Paste: "ctrl+v",

		EditNode:    "enter",
		RemoveNode:  "backspace",
		TriggerNode: "=",

		MuteNode:    "m",
		MuteAllNode: "M",

		RootNoteUp:   "`",
		RootNoteDown: "ù",
		ScaleUp:      "£",
		ScaleDown:    "%",
		TempoUp:      "-",
		TempoDown:    ")",

		Configuration:   "f2",
		FitGridToWindow: "f10",

		Cancel: "esc",

		Help: "?",
		Quit: "ctrl+q",
	}
}

// NewDefaultQwertyKeyMap returns a new default KeyMap for qwerty keyboards.
func NewDefaultQwertyKeyMap() KeyMap {
	return KeyMap{
		Play: " ",

		Up:    "up",
		Right: "right",
		Down:  "down",
		Left:  "left",

		SelectionUp:    "shift+up",
		SelectionRight: "shift+right",
		SelectionDown:  "shift+down",
		SelectionLeft:  "shift+left",

		EditUp:    "ctrl+up",
		EditRight: "ctrl+right",
		EditDown:  "ctrl+down",
		EditLeft:  "ctrl+left",

		EditInput: ".",

		Bank: "tab",

		AddBang:   "1",
		AddEuclid: "2",
		AddPass:   "3",
		AddSpread: "4",
		AddCycle:  "5",
		AddDice:   "6",
		AddToll:   "7",
		AddZone:   "8",
		AddHole:   "9",

		Copy:  "ctrl+c",
		Cut:   "ctrl+x",
		Paste: "ctrl+v",

		EditNode:    "enter",
		RemoveNode:  "backspace",
		TriggerNode: "/",

		MuteNode:    "m",
		MuteAllNode: "M",

		RootNoteUp:   "'",
		RootNoteDown: ";",
		ScaleUp:      "\"",
		ScaleDown:    ":",
		TempoUp:      "=",
		TempoDown:    "-",

		Configuration:   "f2",
		FitGridToWindow: "f10",

		Cancel: "esc",

		Help: "?",
		Quit: "ctrl+q",
	}
}

// NewDefaultQwertyMacKeyMap returns a new default KeyMap for qwerty mac
// keyboards.
func NewDefaultQwertyMacKeyMap() KeyMap {
	return KeyMap{
		Play: " ",

		Up:    "up",
		Right: "right",
		Down:  "down",
		Left:  "left",

		SelectionUp:    "shift+up",
		SelectionRight: "shift+right",
		SelectionDown:  "shift+down",
		SelectionLeft:  "shift+left",

		EditUp:    "ctrl+up",
		EditRight: "ctrl+right",
		EditDown:  "ctrl+down",
		EditLeft:  "ctrl+left",

		EditInput: ".",

		Bank: "tab",

		AddBang:   "1",
		AddEuclid: "2",
		AddPass:   "3",
		AddSpread: "4",
		AddCycle:  "5",
		AddDice:   "6",
		AddToll:   "7",
		AddZone:   "8",
		AddHole:   "9",

		Copy:  "ctrl+c",
		Cut:   "ctrl+x",
		Paste: "ctrl+v",

		EditNode:    "enter",
		RemoveNode:  "backspace",
		TriggerNode: "/",

		MuteNode:    "m",
		MuteAllNode: "M",

		RootNoteUp:   "'",
		RootNoteDown: ";",
		ScaleUp:      "\"",
		ScaleDown:    ":",
		TempoUp:      "=",
		TempoDown:    "-",

		Configuration:   "f2",
		FitGridToWindow: "f10",

		Cancel: "esc",

		Help: "?",
		Quit: "ctrl+q",
	}
}
