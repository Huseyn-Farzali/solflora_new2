package dao

import (
	"Solflora/db"
	"Solflora/state"
)

type TemperatureEntity struct {
	PresentValue     float64
	ControllerOutput float64
	SetPoint         float64
}

type HumidityEntity struct {
	PresentValue float64
}

type MoistureEntity struct {
	PresentValue float64
}

type TuneProfileEntity struct {
	ProportionalGain float64
	IntegralGain     float64
	DerivativeGain   float64
}

func BuildTemperature(modelStateMap map[state.ConditionVariable]float64) TemperatureEntity {
	return TemperatureEntity{
		PresentValue:     modelStateMap[state.TemperaturePV],
		ControllerOutput: modelStateMap[state.TemperatureCO],
		SetPoint:         modelStateMap[state.TemperatureSP],
	}
}

func (tempEntity *TemperatureEntity) Commit() error {
	_, err := db.DB.Exec(`
		INSERT INTO temperature (present_value, controller_output, set_point)
		VALUES ($1, $2, $3)
	`, tempEntity.PresentValue, tempEntity.ControllerOutput, tempEntity.SetPoint)
	return err
}

func BuildHumidity(modelStateMap map[state.ConditionVariable]float64) HumidityEntity {
	return HumidityEntity{
		PresentValue: modelStateMap[state.HumidityPV],
	}
}

func (humEntity *HumidityEntity) Commit() error {
	_, err := db.DB.Exec(`
		INSERT INTO humidity (present_value)
		VALUES ($1)
	`, humEntity.PresentValue)
	return err
}

func BuildMoisture(modelStateMap map[state.ConditionVariable]float64) MoistureEntity {
	return MoistureEntity{
		PresentValue: modelStateMap[state.MoisturePV],
	}
}

func (moistEntity *MoistureEntity) Commit() error {
	_, err := db.DB.Exec(`
		INSERT INTO moisture (present_value)
		VALUES ($1)
	`, moistEntity.PresentValue)
	return err
}

func BuildTuneProfile(tuneStateMap map[state.TuneVariable]float64) TuneProfileEntity {
	return TuneProfileEntity{
		ProportionalGain: tuneStateMap[state.TemperatureKp],
		IntegralGain:     tuneStateMap[state.TemperatureKi],
		DerivativeGain:   tuneStateMap[state.TemperatureKd],
	}
}

func (tuneEntity *TuneProfileEntity) Commit() error {
	_, err := db.DB.Exec(`
		INSERT INTO tune_profile (proportional_gain, integral_gain, derivative_gain)
		VALUES ($1, $2, $3)
	`, tuneEntity.ProportionalGain, tuneEntity.IntegralGain, tuneEntity.DerivativeGain)
	return err
}
