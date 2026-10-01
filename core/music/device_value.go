package music

import (
	"signls/midi"
)

type DeviceValue struct {
	Device     midi.OutputDevice
	GridDevice *midi.OutputDevice
	Enabled    bool
}

func (d DeviceValue) Get() midi.MidiOutput {
	if d.Enabled {
		return d.Device.ID
	}
	return d.GridDevice.ID
}

func (d DeviceValue) Name() string {
	if d.Enabled {
		return d.Device.Name
	}
	return ""
}
