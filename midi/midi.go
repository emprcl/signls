// Package midi provides ways to interact with music/audio midi devices and
// softwares.
package midi

import (
	"errors"
	"log"
	"sync"

	gomidi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

const (
	// Each midi device can receive notes through a dedicated buffered chan.
	// 16 tracks with all steps activated sending notes to the same device
	// at high tempo can results to a lot of midi messages.
	midiBufferSize = 1024
)

// Midi provides a way to interct with midi devices.
type Midi interface {
	Devices() gomidi.OutPorts
	ActiveDevice() drivers.Out
	SetActiveDevice(device int)
	CycleMidiDevices()
	NoteOn(channel uint8, note uint8, velocity uint8)
	NoteOff(channel uint8, note uint8)
	Silence(channel uint8)
	SilenceAll()
	ControlChange(channel, controller, value uint8)
	ProgramChange(channel uint8, value uint8)
	Pitchbend(channel uint8, value int16)
	AfterTouch(channel uint8, value uint8)
	SendClock()
	Close()
}

// midi contains the midi devices state. We use the gomidi package
// for communicating with available devices.
type midi struct {
	// devices holds all the midi devices outputs that are returned by gomidi.
	devices gomidi.OutPorts

	// active holds the active device
	active int

	// Because we want to allow the usage of multiple midi devices at the same
	// time, we start a goroutine for each device that can receive note trigs.
	// The wait group is used when closing the midi devices (waits for all
	// device goroutines to end).
	// The done chan is used to send the end signal to the goroutines.
	// The output chans receives actual midi messages for each devices.
	waitGroup *sync.WaitGroup
	done      chan struct{}
	outputs   []chan gomidi.Message
}

// New creates a new midi. It retrieves the connected midi
// devices and starts a new goroutine for each of them.
func New() (Midi, error) {
	devices := gomidi.GetOutPorts()
	if len(devices) == 0 {
		return nil, errors.New("no midi drivers")
	}
	midi := &midi{
		devices: devices,
	}
	midi.start()
	return midi, nil
}

// Note retruns the string representation of a note
func Note(note uint8) string {
	return gomidi.Note(note).String()
}

func (m *midi) start() {
	var wg sync.WaitGroup
	wg.Add(len(m.devices))
	m.done = make(chan struct{}, len(m.devices))
	for i, device := range m.devices {
		m.outputs = append(m.outputs, make(chan gomidi.Message, midiBufferSize))
		go func(device drivers.Out, done <-chan struct{}, output <-chan gomidi.Message) {
			defer wg.Done()
			send, err := gomidi.SendTo(device)
			if err != nil {
				log.Fatal(err)
			}
			for {
				select {
				case <-done:
					// Before terminating the goroutine, we drain all the
					// remaining messages, ensuring that all the note off
					// signals will be sent before exiting.
					for len(output) > 0 {
						err := send(<-output)
						if err != nil {
							log.Println(err)
						}
					}
					return
				case msg := <-output:
					err := send(msg)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}(device, m.done, m.outputs[i])
	}
	m.waitGroup = &wg
}

// Devices returns all out ports.
func (m *midi) Devices() gomidi.OutPorts {
	return m.devices
}

// CycleMidiDevices cycle through all devices and
// and sets the next one as active
func (m *midi) CycleMidiDevices() {
	if len(m.Devices()) < m.active+2 {
		m.active = 0
		return
	}
	m.active++
}

// SetActiveDevice sets the active device
func (m *midi) SetActiveDevice(device int) {
	if len(m.devices) < device+1 {
		m.active = 0
		return
	}
	m.active = device
}

// ActiveDevice returns the active device
func (m *midi) ActiveDevice() drivers.Out {
	if len(m.devices) < m.active+1 {
		return nil
	}
	return m.devices[m.active]
}

// NoteOn sends a Note On midi meessage to the active device.
func (m *midi) NoteOn(channel uint8, note uint8, velocity uint8) {
	m.outputs[m.active] <- gomidi.NoteOn(channel, note, velocity)
}

// NoteOff sends a Note Off midi meessage to the active device.
func (m *midi) NoteOff(channel uint8, note uint8) {
	m.outputs[m.active] <- gomidi.NoteOff(channel, note)
}

// Silence sends a note off message for every running note on given channel.
func (m *midi) Silence(channel uint8) {
	for _, msg := range gomidi.SilenceChannel(int8(channel)) {
		m.outputs[m.active] <- msg
	}
}

// SilenceAll sends a note off message for every running note on every channel.
func (m *midi) SilenceAll() {
	for c := 0; c < 16; c++ {
		m.Silence(uint8(c))
	}
}

// ControlChange sends a Control Change messages to the active device.
func (m *midi) ControlChange(channel, controller, value uint8) {
	m.outputs[m.active] <- gomidi.ControlChange(channel, controller, value)
}

// ProgramChange sends a Program Change messages to the active device.
func (m *midi) ProgramChange(channel uint8, value uint8) {
	m.outputs[m.active] <- gomidi.ProgramChange(channel, value)
}

// Pitchbend sends a Pitch Bend messages to the active device.
func (m *midi) Pitchbend(channel uint8, value int16) {
	m.outputs[m.active] <- gomidi.Pitchbend(channel, value)
}

// AfterTouch sends a After Touch messages to the active device.
func (m *midi) AfterTouch(channel uint8, value uint8) {
	m.outputs[m.active] <- gomidi.AfterTouch(channel, value)
}

// SendClock sends a Clock midi meessage to the active device.
func (m *midi) SendClock() {
	m.outputs[m.active] <- gomidi.TimingClock()
}

// Close terminates all the device goroutines gracefully.
func (m *midi) Close() {
	defer gomidi.CloseDriver()
	if m.waitGroup == nil {
		return
	}
	for range m.devices {
		m.done <- struct{}{}
	}
	m.waitGroup.Wait()
}
