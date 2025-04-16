package web

import (
	"Solflora/logger"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func WaterPumpControl(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.WaterPumpControl")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.WaterPumpControl | method not allowed: %s", r.Method)
			return
		}

		waterPumpOnStateDuration, err := time.ParseDuration(os.Getenv("WATER_PUMP_ON_STATE_DURATION"))
		if err != nil {
			log.Warn("[WARN] api.web.WaterPumpControl | $env:{WATER_PUMP_ON_STATE_DURATION} is not valid duration – defaulting to 4s")
			waterPumpOnStateDuration = 4 * time.Second
		}

		service.ActivateWaterPump(waterPumpOnStateDuration)
		log.Info("[END] api.web.WaterPumpControl")
	}
}

func AirFanControl(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.AirFanControl")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.AirFanControl | method not allowed: %s", r.Method)
			return
		}

		query := r.URL.Query()
		airFanState, err := mapQueryParamToBoolState(query.Get("state"))
		if err != nil {
			http.Error(w, "Bad Request – state query parameter are not valid", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.AirFanControl | query parameter are not valid | error: %s", err)
			return
		}

		service.UpdateAirFanState(airFanState)
		log.Info("[END] api.web.AirFanControl")
	}
}

func mapQueryParamToBoolState(state string) (bool, error) {
	switch strings.ToLower(state) {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return false, fmt.Errorf("map query param [%s] is not valid", state)
	}
}
