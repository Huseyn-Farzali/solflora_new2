package esp

import (
	"Solflora/logger"
	"encoding/json"
	"net/http"
)

type RequestBody struct {
	TemperaturePV float64 `json:"temp_pv"`
	TemperatureCO float64 `json:"temp_co"`
	TemperatureSP float64 `json:"temp_sp"`
	MoisturePV    float64 `json:"moist_pv"`
	HumidityPV    float64 `json:"humidity_pv"`
}

type ResponseBody struct {
	TemperatureSP    float64 `json:"temp_sp"`
	TemperatureKp    float64 `json:"temp_kp"`
	TemperatureKi    float64 `json:"temp_ki"`
	TemperatureKd    float64 `json:"temp_kd"`
	FanControl       int16   `json:"fan_control"`
	WaterPumpControl int16   `json:"water_pump_control"`
}

func ControlSampler(service *ControlSamplingService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.esp.ControlSampler")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.esp.ControlSampler | method not allowed: %s", r.Method)
			return
		}

		var reqBody RequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Error("[ERROR] api.esp.ControlSampler | invalid request body")
			return
		}
		log.Debugf("[DEBUG] api.esp.ControlSampler() | request body: %+v\n", reqBody)

		respBody, err := service.HandleControlSampling(reqBody)
		if err != nil {
			log.Errorf("[ERROR] api.esp.ControlSampler | handle control sampling failed: %s", err.Error())
			http.Error(w, "Handling control-sampling failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.esp.ControlSampler")
	}

}
