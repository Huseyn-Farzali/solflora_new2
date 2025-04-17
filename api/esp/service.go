package esp

import (
	"Solflora/dao"
	"Solflora/logger"
	"Solflora/state"
)

type ControlSamplingService struct {
	modelState  *state.ModelState
	deviceState *state.DeviceState
	tuneState   *state.TuneState
}

func NewControlSamplingService(
	modelState *state.ModelState,
	deviceState *state.DeviceState,
	tuneState *state.TuneState) *ControlSamplingService {
	return &ControlSamplingService{
		modelState:  modelState,
		deviceState: deviceState,
		tuneState:   tuneState}
}

func (s *ControlSamplingService) HandleControlSampling(req RequestBody) (ResponseBody, error) {
	var log = logger.Logger()
	log.Infof("[START] api.esp.HandleControlSampling")

	var newTemperatureEntity = buildTemperatureEntity(req)
	var newHumidityEntity = buildHumidityEntity(req)
	var newMoistureEntity = buildMoistureEntity(req)

	err := newTemperatureEntity.Commit()
	if err != nil {
		log.Errorf("[ERROR] api.esp.HandlerControlSampling() | temp-entity cannot be committed | %s\n", err.Error())
		return ResponseBody{}, err
	}
	log.Debugf("[DEBUG] api.esp.HandleControlSampling | generated temperature entity: %+v\n", newTemperatureEntity)

	err = newHumidityEntity.Commit()
	if err != nil {
		log.Errorf("[ERROR] api.esp.HandlerControlSampling() | humidity-entity cannot be committed | %s\n", err.Error())
		return ResponseBody{}, err
	}
	log.Debugf("[DEBUG] api.esp.HandleControlSampling | generated humidity entity: %+v\n", newHumidityEntity)

	err = newMoistureEntity.Commit()
	if err != nil {
		log.Errorf("[ERROR] api.esp.HandlerControlSampling() | moist-entity cannot be committed | %s\n", err.Error())
		return ResponseBody{}, err
	}
	log.Debugf("[DEBUG] api.esp.HandleControlSampling | generated moisture entity: %+v\n", newMoistureEntity)

	modelStateMap := s.modelState.GetAll()
	deviceStateMap := s.deviceState.GetAll()
	tuneStateMap := s.tuneState.GetAll()
	responseBody := ResponseBody{
		TemperatureSP:    modelStateMap[state.TemperatureSP],
		TemperatureKp:    tuneStateMap[state.TemperatureKp],
		TemperatureKi:    tuneStateMap[state.TemperatureKi],
		TemperatureKd:    tuneStateMap[state.TemperatureKd],
		FanControl:       boolToInt16(deviceStateMap[state.FanControl]),
		WaterPumpControl: boolToInt16(deviceStateMap[state.WaterPumpControl]),
	}

	log.Debugf("[DEBUG] api.esp.HandleControlSampling | generated response to esp: %+v", responseBody)
	log.Infof("[END] api.esp.HandleControlSampling")
	return responseBody, nil
}

func buildTemperatureEntity(req RequestBody) *dao.TemperatureEntity {

	return &dao.TemperatureEntity{
		PresentValue:     req.TemperaturePV,
		ControllerOutput: req.TemperatureCO,
		SetPoint:         req.TemperatureSP,
	}
}

func buildHumidityEntity(req RequestBody) *dao.HumidityEntity {
	return &dao.HumidityEntity{
		PresentValue: req.HumidityPV,
	}
}

func buildMoistureEntity(req RequestBody) *dao.MoistureEntity {
	return &dao.MoistureEntity{
		PresentValue: req.MoisturePV,
	}
}

func boolToInt16(b bool) int16 {
	if b {
		return 1
	}
	return 0
}
