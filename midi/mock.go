package midi

import (
	gomidi "gitlab.com/gomidi/midi/v2"
)

type Mock struct{}

func (m *Mock) Devices() gomidi.OutPorts                                            { return nil }
func (m *Mock) NoteOn(device MidiOutput, channel uint8, note uint8, velocity uint8) {}
func (m *Mock) NoteOff(device MidiOutput, channel uint8, note uint8)                {}
func (m *Mock) Silence(device MidiOutput, channel uint8)                            {}
func (m *Mock) SilenceAll()                                                         {}
func (m *Mock) ControlChange(device MidiOutput, channel, controller, value uint8)   {}
func (m *Mock) ProgramChange(device MidiOutput, channel uint8, value uint8)         {}
func (m *Mock) Pitchbend(device MidiOutput, channel uint8, value int16)             {}
func (m *Mock) AfterTouch(device MidiOutput, channel uint8, value uint8)            {}
func (m *Mock) SendClock(device MidiOutput)                                         {}
func (m *Mock) TransportStart(device MidiOutput)                                    {}
func (m *Mock) TransportStop(device MidiOutput)                                     {}
func (m *Mock) Send(device MidiOutput, message gomidi.Message)                      {}
func (m *Mock) NewOutDevice(device, fallback string) OutputDevice                   { return OutputDevice{midi: m} }
func (m *Mock) GetOutDevice(device MidiOutput) OutputDevice                         { return OutputDevice{midi: m} }
func (m *Mock) NewInDevice(device, fallback string) InputDevice                     { return InputDevice{midi: m} }
func (m *Mock) GetInDevice(device MidiInput) InputDevice                            { return InputDevice{midi: m} }
func (m *Mock) Listen(cb MidiInCallback) error                                      { return nil }
func (m *Mock) Close()                                                              {}
