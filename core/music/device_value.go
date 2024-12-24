package music

import (
	"signls/midi"
)

type DeviceValue struct {
	Device     midi.Device
	GridDevice *midi.Device
	Enabled    bool
}

func (d DeviceValue) Get() int {
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
