// Package midi provides ways to interact with music/audio midi devices and
// softwares.
package midi

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"

	gomidi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	rtmidi "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

const (
	defaultDevice int = 0

	// Each midi device can receive notes through a dedicated buffered chan.
	// 16 tracks with all steps activated sending notes to the same device
	// at high tempo can results to a lot of midi messages.
	midiBufferSize = 1024
)

// Midi provides a way to interct with midi devices.
type Midi interface {
	Devices() gomidi.OutPorts
	NoteOn(device int, channel uint8, note uint8, velocity uint8)
	NoteOff(device int, channel uint8, note uint8)
	Silence(device int, channel uint8)
	SilenceAll()
	ControlChange(device int, channel, controller, value uint8)
	ProgramChange(device int, channel uint8, value uint8)
	Pitchbend(device int, channel uint8, value int16)
	AfterTouch(device int, channel uint8, value uint8)
	SendClock(device int)
	TransportStart(device int)
	TransportStop(device int)
	NewDevice(device, fallback string) Device
	GetDevice(device int) Device
	Close()
}

// Device represents a midi device.
type Device struct {
	Name     string
	ID       int
	Fallback bool
}

func (d Device) Enabled() bool {
	return d.Name != ""
}

// midi contains the midi devices state. We use the gomidi package
// for communicating with available devices.
type midi struct {
	// devices holds all the midi devices outputs that are returned by gomidi.
	devices gomidi.OutPorts

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
	var m *midi
	if runtime.GOOS != "windows" {
		virtualDevice, err := drivers.Get().(*rtmidi.Driver).OpenVirtualOut("Signls Default Midi Output")
		if err != nil {
			return nil, err
		}
		devices = append(devices, virtualDevice)
	}
	if len(devices) == 0 {
		return nil, errors.New("no midi devices available")
	}
	m = &midi{
		devices: devices,
	}
	m.start()
	return m, nil
}

// Note retruns the string representation of a note
func Note(note uint8) string {
	return gomidi.Note(note).String()
}

// CC retruns the string representation of a cc message
func CC(controller uint8) string {
	return gomidi.ControlChangeName[controller]
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

// NewDevice creates a new device.
func (m *midi) NewDevice(device, fallback string) Device {
	id, err := m.findDeviceIndex(device)
	if err == nil {
		return Device{
			Name: device,
			ID:   id,
		}
	}
	id, err = m.findDeviceIndex(fallback)
	if err == nil {
		return Device{
			Name:     device,
			ID:       id,
			Fallback: true,
		}
	}
	return Device{
		Name:     device,
		ID:       defaultDevice,
		Fallback: true,
	}
}

// Devices returns all out ports.
func (m *midi) Devices() gomidi.OutPorts {
	return m.devices
}

// NoteOn sends a Note On midi meessage to the active device.
func (m *midi) NoteOn(device int, channel uint8, note uint8, velocity uint8) {
	m.outputs[device] <- gomidi.NoteOn(channel, note, velocity)
}

// NoteOff sends a Note Off midi meessage to the active device.
func (m *midi) NoteOff(device int, channel uint8, note uint8) {
	m.outputs[device] <- gomidi.NoteOff(channel, note)
}

// Silence sends a note off message for every running note on given channel.
func (m *midi) Silence(device int, channel uint8) {
	for _, msg := range gomidi.SilenceChannel(int8(channel)) {
		m.outputs[device] <- msg
	}
}

// SilenceAll sends a note off message for every running note on every channel.
func (m *midi) SilenceAll() {
	for device := range m.devices {
		for c := 0; c < 16; c++ {
			m.Silence(device, uint8(c))
		}
	}
}

// ControlChange sends a Control Change messages to the active device.
func (m *midi) ControlChange(device int, channel, controller, value uint8) {
	m.outputs[device] <- gomidi.ControlChange(channel, controller, value)
}

// ProgramChange sends a Program Change messages to the active device.
func (m *midi) ProgramChange(device int, channel uint8, value uint8) {
	m.outputs[device] <- gomidi.ProgramChange(channel, value)
}

// Pitchbend sends a Pitch Bend messages to the active device.
func (m *midi) Pitchbend(device int, channel uint8, value int16) {
	m.outputs[device] <- gomidi.Pitchbend(channel, value)
}

// AfterTouch sends a After Touch messages to the active device.
func (m *midi) AfterTouch(device int, channel uint8, value uint8) {
	m.outputs[device] <- gomidi.AfterTouch(channel, value)
}

// SendClock sends a Clock midi meessage to the active device.
func (m *midi) SendClock(device int) {
	m.outputs[device] <- gomidi.TimingClock()
}

// TransportStart sends a Start midi meessage to the active device.
func (m *midi) TransportStart(device int) {
	m.outputs[device] <- gomidi.Start()
}

// TransportStop sends a Stop midi meessage to the active device.
func (m *midi) TransportStop(device int) {
	m.outputs[device] <- gomidi.Stop()
}

// findDeviceIndex check if the given device is connected
// or fallback on the given fallback device.
func (m *midi) findDeviceIndex(device string) (int, error) {
	for i, d := range m.devices {
		if d.String() == device {
			return i, nil
		}
	}
	return 0, fmt.Errorf("device %s not connected", device)
}

// GetDevice get a midi device per index.
func (m *midi) GetDevice(device int) Device {
	if len(m.devices)-1 < device {
		return Device{Name: m.devices[0].String()}
	}
	if device < 0 {
		index := len(m.devices) - 1
		return Device{Name: m.devices[index].String(), ID: index}
	}
	return Device{Name: m.devices[device].String(), ID: device}
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
