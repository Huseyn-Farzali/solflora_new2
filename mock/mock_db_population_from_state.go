package mock

import (
	"Solflora/dao"
	"Solflora/state"
	"Solflora/util"
)

func Mock_db_population_from_state() {
	var tuneState = state.NewTuneState()
	var modelState = state.NewModelState()

	for i := 0; i < 10000; i++ {
		mock_populate_states(tuneState, modelState)

		var tempEntity dao.TemperatureEntity
		var humEntity dao.HumidityEntity
		var moistEntity dao.MoistureEntity
		var tuneEntity dao.TuneProfileEntity

		tempEntity, humEntity, moistEntity, tuneEntity = mock_build_db_entries(modelState, tuneState)

		mock_commit_db_entries(&tempEntity, &humEntity, &moistEntity, &tuneEntity)
	}
}

func mock_populate_states(tuneState *state.TuneState, modelState *state.ModelState) {
	tuneState.Set(state.TemperatureKp, util.RandomFloat(0, 100))
	tuneState.Set(state.TemperatureKi, util.RandomFloat(0, 100))
	tuneState.Set(state.TemperatureKd, util.RandomFloat(0, 100))

	modelState.Set(state.TemperaturePV, util.RandomFloat(0, 100))
	modelState.Set(state.TemperatureCO, util.RandomFloat(0, 100))
	modelState.Set(state.TemperatureSP, util.RandomFloat(0, 100))
	modelState.Set(state.MoisturePV, util.RandomFloat(0, 100))
	modelState.Set(state.HumidityPV, util.RandomFloat(0, 100))
}

func mock_build_db_entries(
	modelState *state.ModelState,
	tuneState *state.TuneState) (

	dao.TemperatureEntity,
	dao.HumidityEntity,
	dao.MoistureEntity,
	dao.TuneProfileEntity) {

	resultTemp := dao.BuildTemperature(modelState.GetAll())
	resultHum := dao.BuildHumidity(modelState.GetAll())
	resultMoist := dao.BuildMoisture(modelState.GetAll())
	resultTune := dao.BuildTuneProfile(tuneState.GetAll())

	return resultTemp, resultHum, resultMoist, resultTune
}

func mock_commit_db_entries(
	tempEntity *dao.TemperatureEntity,
	humEntity *dao.HumidityEntity,
	moistEntity *dao.MoistureEntity,
	tuneEntity *dao.TuneProfileEntity) {

	tempEntity.Commit()
	humEntity.Commit()
	moistEntity.Commit()
	tuneEntity.Commit()
}
