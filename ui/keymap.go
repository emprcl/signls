package ui

import (
	"cykl/filesystem"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Play key.Binding

	Up    key.Binding
	Right key.Binding
	Down  key.Binding
	Left  key.Binding

	SelectionUp    key.Binding
	SelectionRight key.Binding
	SelectionDown  key.Binding
	SelectionLeft  key.Binding

	EditUp    key.Binding
	EditRight key.Binding
	EditDown  key.Binding
	EditLeft  key.Binding

	AddBang   key.Binding
	AddSpread key.Binding
	AddCycle  key.Binding
	AddDice   key.Binding
	AddQuota  key.Binding
	AddEuclid key.Binding
	AddZone   key.Binding
	AddPass   key.Binding
	AddHole   key.Binding

	Copy  key.Binding
	Cut   key.Binding
	Paste key.Binding

	EditNode    key.Binding
	RemoveNode  key.Binding
	TriggerNode key.Binding

	MuteNode    key.Binding
	MuteAllNode key.Binding

	RootNoteUp   key.Binding
	RootNoteDown key.Binding
	ScaleUp      key.Binding
	ScaleDown    key.Binding
	TempoUp      key.Binding
	TempoDown    key.Binding

	SelectMidiDevice key.Binding

	Cancel key.Binding

	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.AddBang, k.AddSpread, k.AddCycle, k.AddDice, k.AddQuota, k.AddEuclid, k.AddZone, k.AddPass, k.AddHole},
		{k.Up, k.Right, k.Down, k.Left, k.SelectionUp, k.SelectionRight, k.SelectionDown, k.SelectionLeft, k.EditUp, k.EditRight, k.EditDown, k.EditLeft},
		{k.Play, k.EditNode, k.RemoveNode, k.TriggerNode, k.MuteNode, k.MuteAllNode, k.Copy, k.Cut, k.Paste},
		{k.RootNoteUp, k.RootNoteDown, k.ScaleUp, k.ScaleDown, k.Cancel, k.Help, k.Quit},
	}
}

// Direction returns the direction for a given key msg.
func (k keyMap) Direction(msg tea.KeyMsg) string {
	switch {
	case key.Matches(msg, k.Up, k.SelectionUp, k.EditUp):
		return "up"
	case key.Matches(msg, k.Right, k.SelectionRight, k.EditRight):
		return "right"
	case key.Matches(msg, k.Down, k.SelectionDown, k.EditDown):
		return "down"
	case key.Matches(msg, k.Left, k.SelectionLeft, k.EditLeft):
		return "left"
	default:
		return ""
	}
}

// EmitterSymbol returns the emitter symbol from a key msg.
func (k keyMap) EmitterSymbol(msg tea.KeyMsg) string {
	switch {
	case key.Matches(msg, k.AddBang):
		return "b"
	case key.Matches(msg, k.AddCycle):
		return "c"
	case key.Matches(msg, k.AddSpread):
		return "s"
	case key.Matches(msg, k.AddDice):
		return "d"
	case key.Matches(msg, k.AddEuclid):
		return "e"
	case key.Matches(msg, k.AddQuota):
		return "q"
	case key.Matches(msg, k.AddPass):
		return "p"
	case key.Matches(msg, k.AddZone):
		return "z"
	case key.Matches(msg, k.AddHole):
		return "h"
	default:
		return ""
	}
}

// newKeyMap returns the default key mapping.
func newKeyMap(keys filesystem.KeyMap) keyMap {
	return keyMap{
		Play: key.NewBinding(
			key.WithKeys(keys.Play),
			key.WithHelp(keys.Play, "toggle play"),
		),
		Up: key.NewBinding(
			key.WithKeys(keys.Up),
			key.WithHelp(keys.Up, "move cursor up"),
		),
		Right: key.NewBinding(
			key.WithKeys(keys.Right),
			key.WithHelp(keys.Right, "move cursor right | move selected parameter"),
		),
		Down: key.NewBinding(
			key.WithKeys(keys.Down),
			key.WithHelp(keys.Down, "move cursor down"),
		),
		Left: key.NewBinding(
			key.WithKeys(keys.Left),
			key.WithHelp(keys.Left, "move cursor left | move selected parameter"),
		),
		SelectionUp: key.NewBinding(
			key.WithKeys(keys.SelectionUp),
			key.WithHelp(keys.SelectionUp, "move selection up"),
		),
		SelectionRight: key.NewBinding(
			key.WithKeys(keys.SelectionRight),
			key.WithHelp(keys.SelectionRight, "move selection right"),
		),
		SelectionDown: key.NewBinding(
			key.WithKeys(keys.SelectionDown),
			key.WithHelp(keys.SelectionDown, "move selection down"),
		),
		SelectionLeft: key.NewBinding(
			key.WithKeys(keys.SelectionLeft),
			key.WithHelp(keys.SelectionLeft, "move selection left"),
		),
		EditUp: key.NewBinding(
			key.WithKeys(keys.EditUp),
			key.WithHelp(keys.EditUp, "increase selected parameter"),
		),
		EditRight: key.NewBinding(
			key.WithKeys(keys.EditRight),
			key.WithHelp(keys.EditRight, "select parameter mode"),
		),
		EditDown: key.NewBinding(
			key.WithKeys(keys.EditDown),
			key.WithHelp(keys.EditDown, "decrease selected parameter"),
		),
		EditLeft: key.NewBinding(
			key.WithKeys(keys.EditLeft),
			key.WithHelp(keys.EditLeft, "select parameter mode"),
		),
		AddBang: key.NewBinding(
			key.WithKeys(keys.AddBang),
			key.WithHelp(keys.AddBang, "add bang emitter: emits at startup in all directions"),
		),
		AddSpread: key.NewBinding(
			key.WithKeys(keys.AddSpread),
			key.WithHelp(keys.AddSpread, "add spread emitter: emits in all directions"),
		),
		AddCycle: key.NewBinding(
			key.WithKeys(keys.AddCycle),
			key.WithHelp(keys.AddCycle, "add cycle emitter: emits in one direction after another, clockwise"),
		),
		AddDice: key.NewBinding(
			key.WithKeys(keys.AddDice),
			key.WithHelp(keys.AddDice, "add dice emitter: emits in one random direction"),
		),
		AddQuota: key.NewBinding(
			key.WithKeys(keys.AddQuota),
			key.WithHelp(keys.AddQuota, "add quota emitter: emits in all directions after x triggers"),
		),
		AddEuclid: key.NewBinding(
			key.WithKeys(keys.AddEuclid),
			key.WithHelp(keys.AddEuclid, "add euclid emitter: emits in all directions periodically"),
		),
		AddZone: key.NewBinding(
			key.WithKeys(keys.AddZone),
			key.WithHelp(keys.AddZone, "add zone emitter: emits in all directions, and triggers every emitters nearby at the same time."),
		),
		AddPass: key.NewBinding(
			key.WithKeys(keys.AddPass),
			key.WithHelp(keys.AddPass, "add pass emitter: emits in the opposite direction"),
		),
		AddHole: key.NewBinding(
			key.WithKeys(keys.AddHole),
			key.WithHelp(keys.AddHole, "add pass emitter: teleport incoming signal to a specific location"),
		),
		Copy: key.NewBinding(
			key.WithKeys(keys.Copy),
			key.WithHelp(keys.Copy, "copy node"),
		),
		Cut: key.NewBinding(
			key.WithKeys(keys.Cut),
			key.WithHelp(keys.Cut, "cut node"),
		),
		Paste: key.NewBinding(
			key.WithKeys(keys.Paste),
			key.WithHelp(keys.Paste, "paste node"),
		),
		EditNode: key.NewBinding(
			key.WithKeys(keys.EditNode),
			key.WithHelp(keys.EditNode, "edit selected nodes parameters"),
		),
		RemoveNode: key.NewBinding(
			key.WithKeys(keys.RemoveNode),
			key.WithHelp(keys.RemoveNode, "remove selected nodes"),
		),
		TriggerNode: key.NewBinding(
			key.WithKeys(keys.TriggerNode),
			key.WithHelp(keys.TriggerNode, "trigger selected node"),
		),
		MuteNode: key.NewBinding(
			key.WithKeys(keys.MuteNode),
			key.WithHelp(keys.MuteNode, "toggle selected nodes mute"),
		),
		MuteAllNode: key.NewBinding(
			key.WithKeys(keys.MuteAllNode),
			key.WithHelp(keys.MuteAllNode, "mute/unmute all selected nodes"),
		),
		RootNoteUp: key.NewBinding(
			key.WithKeys(keys.RootNoteUp),
			key.WithHelp(keys.RootNoteUp, "increase root note"),
		),
		RootNoteDown: key.NewBinding(
			key.WithKeys(keys.RootNoteDown),
			key.WithHelp(keys.RootNoteDown, "decrease root note"),
		),
		ScaleUp: key.NewBinding(
			key.WithKeys(keys.ScaleUp),
			key.WithHelp(keys.ScaleUp, "increase scale"),
		),
		ScaleDown: key.NewBinding(
			key.WithKeys(keys.ScaleDown),
			key.WithHelp(keys.ScaleDown, "decrease scale"),
		),
		TempoUp: key.NewBinding(
			key.WithKeys(keys.TempoUp),
			key.WithHelp(keys.TempoUp, "increase tempo"),
		),
		TempoDown: key.NewBinding(
			key.WithKeys(keys.TempoDown),
			key.WithHelp(keys.TempoDown, "decrease tempo"),
		),
		SelectMidiDevice: key.NewBinding(
			key.WithKeys(keys.SelectMidiDevice),
			key.WithHelp(keys.SelectMidiDevice, "select midi device"),
		),
		Cancel: key.NewBinding(
			key.WithKeys(keys.Cancel),
			key.WithHelp(keys.Cancel, "cancel selection | exit node parameter edition"),
		),
		Help: key.NewBinding(
			key.WithKeys(keys.Help),
			key.WithHelp(keys.Help, "show help"),
		),
		Quit: key.NewBinding(
			key.WithKeys(keys.Quit),
			key.WithHelp(keys.Quit, "quit"),
		),
	}
}