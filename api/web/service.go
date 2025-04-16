package web

import (
	"Solflora/logger"
	"Solflora/state"
	"context"
	"time"
)

type ControlHandlerService struct {
	deviceState *state.DeviceState
	modelState  *state.ModelState
	tuneState   *state.TuneState
}

func NewControlHandlerService(
	deviceState *state.DeviceState,
	modelState *state.ModelState,
	tuneState *state.TuneState) *ControlHandlerService {
	return &ControlHandlerService{
		deviceState: deviceState,
		modelState:  modelState,
		tuneState:   tuneState}
}

func (s *ControlHandlerService) ActivateWaterPump(duration time.Duration) {
	var log = logger.Logger()

	s.deviceState.Mutex.Lock()
	defer s.deviceState.Mutex.Unlock()

	s.deviceState.ValueMap[state.WaterPumpControl] = true

	if s.deviceState.CancelTimer != nil {
		log.Debug("[DEBUG] api.web.ActivateWaterPump | overriding existing pump-timer")
		s.deviceState.CancelTimer()
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.deviceState.CancelTimer = cancel

	go func() {
		select {
		case <-time.After(duration):
			s.deviceState.Mutex.Lock()
			s.deviceState.ValueMap[state.WaterPumpControl] = false
			s.deviceState.CancelTimer = nil
			s.deviceState.Mutex.Unlock()
		case <-ctx.Done():
		}
	}()
}

func (s *ControlHandlerService) UpdateAirFanState(updatedState bool) {
	var log = logger.Logger()
	s.deviceState.Set(state.FanControl, updatedState)
	log.Debug("[DEBUG] api.web.UpdateAirFanState | updating air-fan to ", updatedState)
}

func (s *ControlHandlerService) UpdateTemperatureSetPoint(updatedSetPoint float64) {
	var log = logger.Logger()
	s.modelState.Set(state.TemperatureSP, updatedSetPoint)
	log.Debug("[DEBUG] api.web.UpdateTemperatureSetPoint | updating temp_sp to ", updatedSetPoint)
}

func (s *ControlHandlerService) ReturnTemperatureControlTuneProfile() TemperatureControlTuneProfileResponseBody {
	var log = logger.Logger()
	tuneStateMap := s.tuneState.GetAll()
	tuneProfile := TemperatureControlTuneProfileResponseBody{
		ProportionalGain: tuneStateMap[state.TemperatureKp],
		IntegralGain:     tuneStateMap[state.TemperatureKi],
		DerivativeGain:   tuneStateMap[state.TemperatureKd],
	}

	log.Debugf("[DEBUG] api.web.ReturnTemperatureControlTuneProfile | returning temp-tune-profile: %+v\n", tuneProfile)
	return tuneProfile
}
