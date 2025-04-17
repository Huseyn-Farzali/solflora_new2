package util

import (
	"Solflora/logger"
	"Solflora/state"
)

func CalculateCO(sp float64, pv float64, integralState *state.TrackingIntegralState, tuneState *state.TuneState) float64 {
	var log = logger.Logger()

	tuneMap := tuneState.GetAll()
	kp := tuneMap[state.TemperatureKp]
	ki := tuneMap[state.TemperatureKi]
	kd := tuneMap[state.TemperatureKd]

	pidErr := sp - pv
	integralState.SetTrackingIntegralValue(integralState.GetTrackingIntegralValue() + pidErr)
	derivative := pidErr - integralState.GetTrackingErrorValue()
	output := kp*pidErr + ki*integralState.GetTrackingErrorValue() + kd*derivative

	integralState.SetTrackingErrorValue(pidErr)
	log.Debugf("[DEBUG] api.esp.ControlSampler() | calculated temp_co: %.4f\n", output)
	return output
}
