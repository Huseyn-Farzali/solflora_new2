package web

import (
	"Solflora/logger"
	"encoding/json"
	"net/http"
	"strconv"
)

type TemperatureControlTuneProfileResponseBody struct {
	ProportionalGain float64 `json:"temp_kp"`
	IntegralGain     float64 `json:"temp_ki"`
	DerivativeGain   float64 `json:"temp_kd"`
}

func TemperatureSetPointControl(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.TemperatureSetPointControl")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.TemperatureSetPointControl | method not allowed: %s", r.Method)
			return
		}

		query := r.URL.Query()
		newTempSP, err := mapQueryParamToF64(query.Get("sp_value"))
		if err != nil {
			http.Error(w, "Bad Request – sp_value query parameter are not valid", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.TemperatureSetPointControl | query parameter are not valid | error: %s", err)
			return
		}

		service.UpdateTemperatureSetPoint(newTempSP)
		log.Info("[END] api.web.TemperatureSetPointControl")
	}
}

func ReturnTemperatureControlTuneProfile(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.ReturnTemperatureControlTuneProfile")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.ReturnTemperatureControlTuneProfile | method not allowed: %s", r.Method)
			return
		}

		respBody := service.ReturnTemperatureControlTuneProfile()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.web.ReturnTemperatureControlTuneProfile")
	}
}

func mapQueryParamToF64(floatS string) (float64, error) {
	return strconv.ParseFloat(floatS, 64)
}
