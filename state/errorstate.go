package state

import "sync"

type TrackingIntegralState struct {
	mutex                 sync.RWMutex
	trackingIntegralValue float64
	trackingErrorValue    float64
}

func NewTrackingIntegralState() *TrackingIntegralState {
	return &TrackingIntegralState{}
}

func (t *TrackingIntegralState) SetTrackingIntegralValue(trackingError float64) {
	t.mutex.Lock()
	t.trackingIntegralValue = trackingError
	t.mutex.Unlock()
}

func (t *TrackingIntegralState) GetTrackingIntegralValue() float64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.trackingIntegralValue
}

func (t *TrackingIntegralState) GetTrackingErrorValue() float64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.trackingErrorValue
}

func (t *TrackingIntegralState) SetTrackingErrorValue(trackingError float64) {
	t.mutex.Lock()
	t.trackingErrorValue = trackingError
	t.mutex.Unlock()
}
