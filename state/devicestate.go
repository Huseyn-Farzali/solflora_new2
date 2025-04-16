package state

import (
	"context"
	"sync"
)

type DeviceControlVariable string

const (
	FanControl       DeviceControlVariable = "fan_control"
	WaterPumpControl DeviceControlVariable = "pump_water"
)

type DeviceState struct {
	Mutex       sync.RWMutex
	ValueMap    map[DeviceControlVariable]bool
	CancelTimer context.CancelFunc
}

func NewDeviceState() *DeviceState {
	return &DeviceState{
		ValueMap: make(map[DeviceControlVariable]bool),
	}
}

func (state *DeviceState) GetAll() map[DeviceControlVariable]bool {
	state.Mutex.RLock()
	defer state.Mutex.RUnlock()

	copyMap := make(map[DeviceControlVariable]bool, len(state.ValueMap))
	for k, v := range state.ValueMap {
		copyMap[k] = v
	}

	return copyMap
}

func (state *DeviceState) Set(variable DeviceControlVariable, value bool) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	state.ValueMap[variable] = value
}
