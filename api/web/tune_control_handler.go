package web

import (
	"Solflora/logger"
	"encoding/json"
	"net/http"
)

type TemperatureControlTuneProfileRequestBody struct {
	ProportionalGain float64 `json:"temp_kp"`
	IntegralGain     float64 `json:"temp_ki"`
	DerivativeGain   float64 `json:"temp_kd"`
}

type TemperatureControlTuneProfileResponseBody struct {
	ProportionalGain float64 `json:"temp_kp"`
	IntegralGain     float64 `json:"temp_ki"`
	DerivativeGain   float64 `json:"temp_kd"`
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

func SetTemperatureControlTuneProfile(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.SetTemperatureControlTuneProfile")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.SetTemperatureControlTuneProfile | method not allowed: %s", r.Method)
			return
		}

		var reqBody TemperatureControlTuneProfileRequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Error("[ERROR] api.web.SetTemperatureControlTuneProfile | invalid request body: %s", err.Error())
			return
		}
		log.Debugf("[DEBUG] api.web.SetTemperatureControlTuneProfile | request body: %+v\n", reqBody)

		err := service.SetTemperatureControlTuneProfile(reqBody)
		if err != nil {
			http.Error(w, "Internal Server Error – failed to commit new tune-profile", http.StatusInternalServerError)
			log.Errorf("[ERROR] api.web.SetTemperatureControlTuneProfile | failed to commit new tune-profile: %s\n", err.Error())
		}

		respBody := TemperatureControlTuneProfileResponseBody{
			ProportionalGain: reqBody.ProportionalGain,
			IntegralGain:     reqBody.IntegralGain,
			DerivativeGain:   reqBody.DerivativeGain,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.SetTemperatureControlTuneProfile")
	}
}
