package state

import "sync"

type ConditionVariable string

const (
	TemperaturePV ConditionVariable = "temp_pv"
	TemperatureCO ConditionVariable = "temp_co"
	TemperatureSP ConditionVariable = "temp_sp"
	MoisturePV    ConditionVariable = "moist_pv"
	HumidityPV    ConditionVariable = "humidity_pv"
)

type ModelState struct {
	mutex    sync.RWMutex
	valueMap map[ConditionVariable]float64
}

func NewModelState() *ModelState {
	return &ModelState{
		valueMap: make(map[ConditionVariable]float64),
	}
}

func (state *ModelState) GetAll() map[ConditionVariable]float64 {
	state.mutex.RLock()
	defer state.mutex.RUnlock()

	copyMap := make(map[ConditionVariable]float64, len(state.valueMap))
	for k, v := range state.valueMap {
		copyMap[k] = v
	}

	return copyMap
}

func (state *ModelState) Set(variable ConditionVariable, value float64) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	state.valueMap[variable] = value
}
