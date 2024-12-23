package midi

import (
	gomidi "gitlab.com/gomidi/midi/v2"
)

type Mock struct{}

func (m *Mock) Devices() gomidi.OutPorts                                     { return nil }
func (m *Mock) NoteOn(device int, channel uint8, note uint8, velocity uint8) {}
func (m *Mock) NoteOff(device int, channel uint8, note uint8)                {}
func (m *Mock) Silence(device int, channel uint8)                            {}
func (m *Mock) SilenceAll(device int)                                        {}
func (m *Mock) ControlChange(device int, channel, controller, value uint8)   {}
func (m *Mock) ProgramChange(device int, channel uint8, value uint8)         {}
func (m *Mock) Pitchbend(device int, channel uint8, value int16)             {}
func (m *Mock) AfterTouch(device int, channel uint8, value uint8)            {}
func (m *Mock) SendClock(device int)                                         {}
func (m *Mock) TransportStart(device int)                                    {}
func (m *Mock) TransportStop(device int)                                     {}
func (m *Mock) NewDevice(device string) Device                               { return Device{} }
func (m *Mock) GetDevice(device int) Device                                  { return Device{} }
func (m *Mock) Close()                                                       {}
