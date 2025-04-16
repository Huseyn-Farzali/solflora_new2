package state

import "sync"

type TuneVariable string

const (
	TemperatureKp TuneVariable = "temp_kp"
	TemperatureKd TuneVariable = "temp_kd"
	TemperatureKi TuneVariable = "temp_ki"
)

type TuneState struct {
	mutex    sync.RWMutex
	valueMap map[TuneVariable]float64
}

func NewTuneState() *TuneState {
	return &TuneState{
		valueMap: make(map[TuneVariable]float64),
	}
}

func (state *TuneState) GetAll() map[TuneVariable]float64 {
	state.mutex.RLock()
	defer state.mutex.RUnlock()

	copyMap := make(map[TuneVariable]float64, len(state.valueMap))
	for k, v := range state.valueMap {
		copyMap[k] = v
	}

	return copyMap
}

func (state *TuneState) Set(variable TuneVariable, value float64) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.valueMap[variable] = value
}
