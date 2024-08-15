package midi

import (
	gomidi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type Mock struct{}

func (m *Mock) Devices() gomidi.OutPorts                         { return nil }
func (m *Mock) ActiveDevice() drivers.Out                        { return nil }
func (m *Mock) SetActiveDevice(device int)                       {}
func (m *Mock) CycleMidiDevices()                                {}
func (m *Mock) NoteOn(channel uint8, note uint8, velocity uint8) {}
func (m *Mock) NoteOff(channel uint8, note uint8)                {}
func (m *Mock) Silence(channel uint8)                            {}
func (m *Mock) ControlChange(channel, controller, value uint8)   {}
func (m *Mock) ProgramChange(channel uint8, value uint8)         {}
func (m *Mock) Pitchbend(channel uint8, value int16)             {}
func (m *Mock) AfterTouch(channel uint8, value uint8)            {}
func (m *Mock) SendClock()                                       {}
func (m *Mock) Close()                                           {}
