package web

import (
	"Solflora/dao"
	"Solflora/db"
	"Solflora/logger"
	"Solflora/state"
	"context"
	"fmt"
	"strings"
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

func (s *ControlHandlerService) ReturnTemperatureSetPoint() float64 {
	var log = logger.Logger()
	tempSp := s.modelState.GetAll()[state.TemperatureSP]
	log.Debug("[DEBUG] api.web.ReturnTemperatureSetPoint | returning temperature set-point: ", tempSp)
	return tempSp
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

func (s *ControlHandlerService) SetTemperatureControlTuneProfile(profile TemperatureControlTuneProfileRequestBody) error {
	var log = logger.Logger()
	s.tuneState.Set(state.TemperatureKp, profile.ProportionalGain)
	s.tuneState.Set(state.TemperatureKi, profile.IntegralGain)
	s.tuneState.Set(state.TemperatureKd, profile.DerivativeGain)

	updatedTuneState := s.tuneState.GetAll()
	log.Debugf("[DEBUG] api.web.SetTemperatureControlTuneProfile | new temp-tune-profile: %+v\n", updatedTuneState)

	newTuneProfileEntry := dao.TuneProfileEntity{
		ProportionalGain: updatedTuneState[state.TemperatureKp],
		IntegralGain:     updatedTuneState[state.TemperatureKi],
		DerivativeGain:   updatedTuneState[state.TemperatureKd],
	}
	err := newTuneProfileEntry.Commit()
	if err != nil {
		log.Errorf("[ERROR] api.web.SetTemperatureControlTuneProfile | failed to commit new profile entity: %s", err.Error())
		return err
	}
	log.Debugf("[DEBUG] api.web.SetTemperatureControlTuneProfile | new temp-tune-profile-entry: %+v\n", newTuneProfileEntry)

	return nil
}

func (s *ControlHandlerService) ReturnMoistureChartData(interval time.Duration, sampling time.Duration) ([]MoistureChartDataEntry, error) {
	var log = logger.Logger()

	if sampling <= 0 {
		return nil, fmt.Errorf("sampling is less than 0")
	}

	pgInterval := toPostgresIntervalString(interval)
	log.Debugf("[DEBUG] api.web.ReturnMoistureChartData | pgInterval: %s", pgInterval)
	rows, err := db.DB.Query(`
		SELECT present_value, created_at
		FROM moisture
		WHERE created_at >= NOW() - INTERVAL '` + pgInterval + `'`)

	if err != nil {
		log.Errorf("[ERROR] api.web.ReturnMoistureChartData | failed to retrieve humidity chart data (int: %s): %s", interval, err.Error())
		return nil, err
	}
	defer rows.Close()

	var fullIntervalRangeData []MoistureChartDataEntry
	for rows.Next() {
		var entry MoistureChartDataEntry
		if err := rows.Scan(&entry.MoisturePV, &entry.Timestamp); err != nil {
			log.Errorf("[ERROR] api.web.ReturnMoistureChartData | failed to scan row: %s", err.Error())
			return nil, err
		}
		fullIntervalRangeData = append(fullIntervalRangeData, entry)
	}

	if len(fullIntervalRangeData) == 0 {
		log.Errorf("[ERROR] api.web.ReturnMoistureChartData | failed to retrieve temperature chart data or chart data is empty (int: %s)", interval)
		return nil, fmt.Errorf("no temperature chart could be fetched")
	}

	targetLength := int(interval.Milliseconds() / sampling.Milliseconds())

	if targetLength <= 0 {
		log.Errorf("[ERROR] api.web.ReturnMoistureChartData | targetCount <= 0 (int: %s, samp: %s)", interval, sampling)
		return nil, fmt.Errorf("targetCount is less than 0")
	}

	var lastTime time.Time
	sampledData := make([]MoistureChartDataEntry, targetLength)
	for _, entry := range fullIntervalRangeData {
		parsedTime, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			continue
		}

		if lastTime.IsZero() || parsedTime.Sub(lastTime) >= sampling {
			sampledData = append(sampledData, entry)
			lastTime = parsedTime
		}
	}

	return sampledData, nil
}

func (s *ControlHandlerService) ReturnHumidityChartData(interval time.Duration, sampling time.Duration) ([]HumidityChartDataEntry, error) {
	var log = logger.Logger()

	if sampling <= 0 {
		return nil, fmt.Errorf("sampling is less than 0")
	}

	pgInterval := toPostgresIntervalString(interval)
	log.Debugf("[DEBUG] api.web.ReturnHumidityChartData | pgInterval: %s", pgInterval)
	rows, err := db.DB.Query(`
		SELECT present_value, created_at
		FROM humidity
		WHERE created_at >= NOW() - INTERVAL '` + pgInterval + `'`)

	if err != nil {
		log.Errorf("[ERROR] api.web.ReturnHumidityChartData | failed to retrieve humidity chart data (int: %s): %s", interval, err.Error())
		return nil, err
	}
	defer rows.Close()

	var fullIntervalRangeData []HumidityChartDataEntry
	for rows.Next() {
		var entry HumidityChartDataEntry
		if err := rows.Scan(&entry.HumidityPV, &entry.Timestamp); err != nil {
			log.Errorf("[ERROR] api.web.ReturnHumidityChartData | failed to scan row: %s", err.Error())
			return nil, err
		}
		fullIntervalRangeData = append(fullIntervalRangeData, entry)
	}

	if len(fullIntervalRangeData) == 0 {
		log.Errorf("[ERROR] api.web.ReturnHumidityChartData | failed to retrieve temperature chart data or chart data is empty (int: %s)", interval)
		return nil, fmt.Errorf("no temperature chart could be fetched")
	}

	targetLength := int(interval.Milliseconds() / sampling.Milliseconds())

	if targetLength <= 0 {
		log.Errorf("[ERROR] api.web.ReturnHumidityChartData | targetCount <= 0 (int: %s, samp: %s)", interval, sampling)
		return nil, fmt.Errorf("targetCount is less than 0")
	}

	var lastTime time.Time
	sampledData := make([]HumidityChartDataEntry, targetLength)
	for _, entry := range fullIntervalRangeData {
		parsedTime, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			continue
		}

		if lastTime.IsZero() || parsedTime.Sub(lastTime) >= sampling {
			sampledData = append(sampledData, entry)
			lastTime = parsedTime
		}
	}

	return sampledData, nil
}

func (s *ControlHandlerService) ReturnTemperatureChartData(interval time.Duration, sampling time.Duration) ([]TemperatureChartDataEntry, error) {
	var log = logger.Logger()

	if sampling <= 0 {
		return nil, fmt.Errorf("sampling is less than 0")
	}

	pgInterval := toPostgresIntervalString(interval)
	log.Debugf("[DEBUG] api.web.ReturnTemperatureChartData | pgInterval: %s", pgInterval)
	rows, err := db.DB.Query(`
		SELECT present_value, controller_output, set_point, created_at
		FROM temperature
		WHERE created_at >= NOW() - INTERVAL '` + pgInterval + `'`)

	if err != nil {
		log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | failed to retrieve temperature chart data (int: %s): %s", interval, err.Error())
		return nil, err
	}
	defer rows.Close()

	var fullIntervalRangeData []TemperatureChartDataEntry
	for rows.Next() {
		var entry TemperatureChartDataEntry
		if err := rows.Scan(&entry.TemperaturePV, &entry.TemperatureCO, &entry.TemperatureSP, &entry.Timestamp); err != nil {
			log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | failed to scan row: %s", err.Error())
			return nil, err
		}
		fullIntervalRangeData = append(fullIntervalRangeData, entry)
	}

	if len(fullIntervalRangeData) == 0 {
		log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | failed to retrieve temperature chart data or chart data is empty (int: %s)", interval)
		return nil, fmt.Errorf("no temperature chart could be fetched")
	}

	targetLength := int(interval.Milliseconds() / sampling.Milliseconds())

	if targetLength <= 0 {
		log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | targetCount <= 0 (int: %s, samp: %s)", interval, sampling)
		return nil, fmt.Errorf("targetCount is less than 0")
	}

	var lastTime time.Time
	sampledData := make([]TemperatureChartDataEntry, targetLength)
	for _, entry := range fullIntervalRangeData {
		parsedTime, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err != nil {
			continue
		}

		if lastTime.IsZero() || parsedTime.Sub(lastTime) >= sampling {
			sampledData = append(sampledData, entry)
			lastTime = parsedTime
		}
	}

	return sampledData, nil
}

func toPostgresIntervalString(d time.Duration) string {
	d = d.Truncate(time.Millisecond)

	hours := int64(d / time.Hour)
	d -= time.Duration(hours) * time.Hour
	minutes := int64(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute
	seconds := int64(d / time.Second)
	d -= time.Duration(seconds) * time.Second
	milliseconds := int64(d / time.Millisecond)

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", seconds))
	}
	if milliseconds > 0 {
		parts = append(parts, fmt.Sprintf("%d milliseconds", milliseconds))
	}

	if len(parts) == 0 {
		return "0 seconds"
	}

	return strings.Join(parts, " ")
}
