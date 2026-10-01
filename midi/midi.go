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
	defaultOutDevice MidiOutput = 0
	defaultInDevice  MidiInput  = -1

	// Each midi device can receive notes through a dedicated buffered chan.
	// 16 tracks with all steps activated sending notes to the same device
	// at high tempo can results to a lot of midi messages.
	midiBufferSize = 1024
)

// There are two classes of midi devices: Inputs and Outputs
type MidiInput int
type MidiOutput int

// Midi provides a way to interct with midi devices.
type Midi interface {
	Devices() gomidi.OutPorts
	NoteOn(device MidiOutput, channel uint8, note uint8, velocity uint8)
	NoteOff(device MidiOutput, channel uint8, note uint8)
	Silence(device MidiOutput, channel uint8)
	SilenceAll()
	ControlChange(device MidiOutput, channel, controller, value uint8)
	ProgramChange(device MidiOutput, channel uint8, value uint8)
	Pitchbend(device MidiOutput, channel uint8, value int16)
	AfterTouch(device MidiOutput, channel uint8, value uint8)
	SendClock(device MidiOutput)
	TransportStart(device MidiOutput)
	TransportStop(device MidiOutput)
	Send(device MidiOutput, message gomidi.Message)
	NewOutDevice(device, fallback string) OutputDevice
	GetOutDevice(device MidiOutput) OutputDevice
	NewInDevice(device, fallback string) InputDevice
	GetInDevice(device MidiInput) InputDevice
	Listen(cb MidiInCallback) error
	Close()
}

// Device represents a midi output device.
type OutputDevice struct {
	midi     Midi
	Name     string
	ID       MidiOutput
	Fallback bool
	// When choosing next and previous devices, this sets whether or not a
	// no-name (i.e. disabled) device is part of the rotation.
	Nullable bool
}

func (d OutputDevice) Enabled() bool {
	return d.Name != ""
}

func (d OutputDevice) Next() OutputDevice {
	next_id := d.ID + 1
	if !d.Enabled() && d.Nullable {
		next_id = 0
	}
	next := d.midi.GetOutDevice(next_id)
	next.Nullable = d.Nullable
	if next.ID != next_id && d.Nullable {
		// Wraparound, insert the disabled null device here
		return OutputDevice{midi: d.midi, ID: defaultOutDevice, Fallback: true, Nullable: true}
	}
	return next
}

func (d OutputDevice) Prev() OutputDevice {
	if d.ID == 0 && d.Enabled() && d.Nullable {
		// Wraparound, insert the disabled null device here
		return OutputDevice{midi: d.midi, ID: defaultOutDevice, Fallback: true, Nullable: true}
	}
	next := d.midi.GetOutDevice(d.ID - 1)
	next.Nullable = d.Nullable
	return next
}

// Device represents a midi input device.
type InputDevice struct {
	midi     Midi
	Name     string
	ID       MidiInput
	Fallback bool
	// NOTE: Input devices are optional, so the Next/Prev sequence always includes
	// a no-name (i.e. disabled) device.
}

func (d InputDevice) Enabled() bool {
	return d.Name != ""
}

func (d InputDevice) Prev() InputDevice {
	return d.midi.GetInDevice(MidiInput(d.ID - 1))
}

func (d InputDevice) Next() InputDevice {
	return d.midi.GetInDevice(MidiInput(d.ID + 1))
}

// Message received from a device. The passed message will be one of the
// types defined here
type InMessage struct {
	Midi    Midi
	Device  MidiInput
	Message interface{}
	Raw     gomidi.Message
}

// Received NoteOn message with velocity > 0
type NoteStart struct {
	Channel, Note, Velocity uint8
}

// Received ControlChange message
type ControlChange struct {
	Channel, Controller, Value uint8
}

// Received AfterTouch message
type AfterTouch struct {
	Channel, Pressure uint8
}

// Received ProgramChange message
type ProgramChange struct {
	Channel, Program uint8
}

// Received PitchBend message
type PitchBend struct {
	Channel  uint8
	Relative int16
	Absolute uint16
}

// Midi events are received by callback
type MidiInCallback func(InMessage)

// midi contains the midi devices state. We use the gomidi package
// for communicating with available devices.
type midi struct {
	// these hold all the midi devices that are returned by gomidi.
	outDevices gomidi.OutPorts
	inDevices  gomidi.InPorts

	// Because we want to allow the usage of multiple midi devices at the same
	// time, we start a goroutine for each device that can receive note trigs.
	// The wait group is used when closing the midi devices (waits for all
	// device goroutines to end).
	// The done chan is used to send the end signal to the goroutines.
	// The output chans receives actual midi messages for each devices.
	waitGroup *sync.WaitGroup
	done      chan struct{}
	outputs   []chan gomidi.Message

	// When listening, we listen to all devices, so there may be many
	// concurrent listeners
	stop []func()
}

// New creates a new midi. It retrieves the connected midi
// devices and starts a new goroutine for each of them.
func New() (Midi, error) {
	outDevices := gomidi.GetOutPorts()
	inDevices := gomidi.GetInPorts()
	var m *midi
	if runtime.GOOS != "windows" {
		virtualDevice, err := drivers.Get().(*rtmidi.Driver).OpenVirtualOut("Signls Default Midi Output")
		if err != nil {
			return nil, err
		}
		outDevices = append(outDevices, virtualDevice)
	}
	if len(outDevices) == 0 {
		return nil, errors.New("no midi devices available")
	}
	// NOTE: It is not a problem if there are no input devices available
	m = &midi{
		outDevices: outDevices,
		inDevices:  inDevices,
		stop:       make([]func(), 0),
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
	// We launch a goroutine to handle output ports
	var wg sync.WaitGroup
	wg.Add(len(m.outDevices))
	m.done = make(chan struct{}, len(m.outDevices))
	for i, device := range m.outDevices {
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
	// For input devices, these are subscribed opportunistically
}

// NewDevice creates a new device.
func (m *midi) NewOutDevice(device, fallback string) OutputDevice {
	id, err := m.findOutDeviceIndex(device)
	if err == nil {
		return OutputDevice{
			midi: m,
			Name: device,
			ID:   id,
		}
	}
	id, err = m.findOutDeviceIndex(fallback)
	if err == nil {
		return OutputDevice{
			midi:     m,
			Name:     device,
			ID:       id,
			Fallback: true,
		}
	}
	return OutputDevice{
		midi:     m,
		Name:     device,
		ID:       defaultOutDevice,
		Fallback: true,
	}
}

// Devices returns all out ports.
func (m *midi) Devices() gomidi.OutPorts {
	return m.outDevices
}

// NoteOn sends a Note On midi meessage to the active device.
func (m *midi) NoteOn(device MidiOutput, channel uint8, note uint8, velocity uint8) {
	m.outputs[device] <- gomidi.NoteOn(channel, note, velocity)
}

// NoteOff sends a Note Off midi meessage to the active device.
func (m *midi) NoteOff(device MidiOutput, channel uint8, note uint8) {
	m.outputs[device] <- gomidi.NoteOff(channel, note)
}

// Silence sends a note off message for every running note on given channel.
func (m *midi) Silence(device MidiOutput, channel uint8) {
	for _, msg := range gomidi.SilenceChannel(int8(channel)) {
		m.outputs[device] <- msg
	}
}

// SilenceAll sends a note off message for every running note on every channel.
func (m *midi) SilenceAll() {
	for device := range m.outDevices {
		for c := 0; c < 16; c++ {
			m.Silence(MidiOutput(device), uint8(c))
		}
	}
}

// ControlChange sends a Control Change messages to the active device.
func (m *midi) ControlChange(device MidiOutput, channel, controller, value uint8) {
	m.outputs[device] <- gomidi.ControlChange(channel, controller, value)
}

// ProgramChange sends a Program Change messages to the active device.
func (m *midi) ProgramChange(device MidiOutput, channel uint8, value uint8) {
	m.outputs[device] <- gomidi.ProgramChange(channel, value)
}

// Pitchbend sends a Pitch Bend messages to the active device.
func (m *midi) Pitchbend(device MidiOutput, channel uint8, value int16) {
	m.outputs[device] <- gomidi.Pitchbend(channel, value)
}

// AfterTouch sends a After Touch messages to the active device.
func (m *midi) AfterTouch(device MidiOutput, channel uint8, value uint8) {
	m.outputs[device] <- gomidi.AfterTouch(channel, value)
}

// SendClock sends a Clock midi meessage to the active device.
func (m *midi) SendClock(device MidiOutput) {
	m.outputs[device] <- gomidi.TimingClock()
}

// TransportStart sends a Start midi meessage to the active device.
func (m *midi) TransportStart(device MidiOutput) {
	m.outputs[device] <- gomidi.Start()
}

// TransportStop sends a Stop midi meessage to the active device.
func (m *midi) TransportStop(device MidiOutput) {
	m.outputs[device] <- gomidi.Stop()
}

// Send sends an arbitrary message
func (m *midi) Send(device MidiOutput, msg gomidi.Message) {
	m.outputs[device] <- msg
}

// findDeviceIndex check if the given device is connected
// or fallback on the given fallback device.
func (m *midi) findOutDeviceIndex(device string) (MidiOutput, error) {
	for i, d := range m.outDevices {
		if d.String() == device {
			return MidiOutput(i), nil
		}
	}
	return MidiOutput(0), fmt.Errorf("device %s not connected", device)
}

// GetOutDevice get a midi output device per index.
func (m *midi) GetOutDevice(device MidiOutput) OutputDevice {
	if len(m.outDevices)-1 < int(device) {
		return m.GetOutDevice(MidiOutput(0))
	}
	if device < 0 {
		return m.GetOutDevice(MidiOutput(len(m.outDevices) - 1))
	}
	return OutputDevice{midi: m, Name: m.outDevices[device].String(), ID: MidiOutput(device)}
}

// getDisabledOutDevice gets a midi output device that reads as disabled
func (m *midi) getDisabledOutDevice() OutputDevice {
	return OutputDevice{midi: m, Name: "", ID: defaultOutDevice, Fallback: true}
}

func (m *midi) findInDeviceIndex(device string) (MidiInput, error) {
	for i, d := range m.inDevices {
		if d.String() == device {
			return MidiInput(i), nil
		}
	}
	return defaultInDevice, fmt.Errorf("device %s not connected", device)
}

func (m *midi) NewInDevice(device, fallback string) InputDevice {
	id, err := m.findInDeviceIndex(device)
	if err == nil {
		return InputDevice{
			midi: m,
			Name: device,
			ID:   id,
		}
	}
	id, err = m.findInDeviceIndex(fallback)
	if err == nil {
		return InputDevice{
			midi:     m,
			Name:     device,
			ID:       id,
			Fallback: true,
		}
	}
	return InputDevice{
		midi:     m,
		Name:     device,
		ID:       defaultInDevice,
		Fallback: true,
	}
}

func (m *midi) GetInDevice(device MidiInput) InputDevice {
	if device < defaultInDevice-1 && len(m.inDevices) > 0 {
		// Wrap back to top
		return m.GetInDevice(MidiInput(len(m.inDevices) - 1))
	} else if device < 0 || len(m.inDevices)-1 < int(device) {
		// Insert disabled device at wrap points
		return m.getDisabledInDevice()
	}
	return InputDevice{midi: m, Name: m.inDevices[device].String(), ID: device}
}

// getDisabledInDevice returns an midi input device that reads as disabled
func (m *midi) getDisabledInDevice() InputDevice {
	return InputDevice{midi: m, Name: "", ID: defaultInDevice, Fallback: true}
}

// Listen subscribes to messages from all connected input devices
func (m *midi) Listen(cb MidiInCallback) error {
	for i, device := range m.inDevices {
		stop, err := gomidi.ListenTo(device, func(msg gomidi.Message, timestamps int32) {
			var channel, key, velocity, controller, value, pressure, program uint8
			var relative int16
			var absolute uint16
			var message interface{} = nil
			switch {
			case msg.GetNoteStart(&channel, &key, &velocity):
				message = &NoteStart{
					Channel:  channel,
					Note:     key,
					Velocity: velocity,
				}
			case msg.GetControlChange(&channel, &controller, &value):
				message = &ControlChange{
					Channel:    channel,
					Controller: controller,
					Value:      value,
				}
			case msg.GetAfterTouch(&channel, &pressure):
				message = &AfterTouch{
					Channel:  channel,
					Pressure: pressure,
				}
			case msg.GetProgramChange(&channel, &program):
				message = &ProgramChange{
					Channel: channel,
					Program: program,
				}
			case msg.GetPitchBend(&channel, &relative, &absolute):
				message = &PitchBend{
					Channel:  channel,
					Relative: relative,
					Absolute: absolute,
				}
			}
			cb(InMessage{
				Midi:    m,
				Device:  MidiInput(i),
				Message: message,
				Raw:     msg,
			})
		})
		if err == nil {
			m.stop = append(m.stop, stop)
		} else {
			return err
		}
	}
	return nil
}

// Close terminates all the device goroutines gracefully.
func (m *midi) Close() {
	defer gomidi.CloseDriver()
	for _, stop := range m.stop {
		if stop != nil {
			stop()
		}
	}
	if m.waitGroup == nil {
		return
	}
	for range m.outDevices {
		m.done <- struct{}{}
	}
	m.waitGroup.Wait()
}
