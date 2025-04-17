package web

import (
	"Solflora/logger"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type TemperatureSetPointResponseBody struct {
	TemperatureSP float64 `json:"temp_sp"`
}

type TemperatureChartDataEntry struct {
	TemperaturePV float64 `json:"temp_pv"`
	TemperatureCO float64 `json:"temp_co"`
	TemperatureSP float64 `json:"temp_sp"`
	Timestamp     string  `json:"time"`
}

type HumidityChartDataEntry struct {
	HumidityPV float64 `json:"humidity"`
	Timestamp  string  `json:"time"`
}

type MoistureChartDataEntry struct {
	MoisturePV float64 `json:"moisture"`
	Timestamp  string  `json:"time"`
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

func ReturnTemperatureSetPoint(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.TemperatureSetPointControl")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.TemperatureSetPointControl | method not allowed: %s", r.Method)
			return
		}

		respBody := TemperatureSetPointResponseBody{
			TemperatureSP: service.ReturnTemperatureSetPoint(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.web.TemperatureSetPointControl")
	}
}

func ReturnMoistureChartData(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.ReturnMoistureChartData")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.ReturnMoistureChartData | method not allowed: %s", r.Method)
			return
		}

		query := r.URL.Query()
		interval, err := mapQueryParamToDuration(query.Get("interval"))
		if err != nil {
			http.Error(w, "Invalid interval query parameter", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.ReturnMoistureChartData | query parameter are not valid | error: %s", err)
		}
		sampling, err := mapQueryParamToDuration(query.Get("sampling"))
		if err != nil {
			http.Error(w, "Invalid sampling query parameter", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.ReturnMoistureChartData | sampling parameter are not valid | error: %s", err)
		}

		entries, err := service.ReturnMoistureChartData(interval, sampling)
		respBody := mapTimeStampToSpecifiedFormatForMoisture(entries)

		if err != nil {
			http.Error(w, "Internal Server Error – ReturnMoistureChartData failed", http.StatusInternalServerError)
			log.Errorf("[ERROR] api.web.ReturnMoistureChartData | ReturnMoistureChartData failed | err: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.web.ReturnMoistureChartData")
	}
}

func ReturnHumidityChartData(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.ReturnHumidityChartData")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.ReturnHumidityChartData | method not allowed: %s", r.Method)
			return
		}

		query := r.URL.Query()
		interval, err := mapQueryParamToDuration(query.Get("interval"))
		if err != nil {
			http.Error(w, "Invalid interval query parameter", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.ReturnHumidityChartData | query parameter are not valid | error: %s", err)
		}
		sampling, err := mapQueryParamToDuration(query.Get("sampling"))
		if err != nil {
			http.Error(w, "Invalid sampling query parameter", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.ReturnHumidityChartData | sampling parameter are not valid | error: %s", err)
		}

		entries, err := service.ReturnHumidityChartData(interval, sampling)
		respBody := mapTimeStampToSpecifiedFormatForHumidity(entries)

		if err != nil {
			http.Error(w, "Internal Server Error – ReturnTemperatureChartData failed", http.StatusInternalServerError)
			log.Errorf("[ERROR] api.web.ReturnHumidityChartData | ReturnTemperatureChartData failed | err: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.web.ReturnHumidityChartData")
	}
}

func ReturnTemperatureChartData(service *ControlHandlerService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var log = logger.Logger()
		log.Info("[START] api.web.ReturnTemperatureChartData")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | method not allowed: %s", r.Method)
			return
		}

		query := r.URL.Query()
		interval, err := mapQueryParamToDuration(query.Get("interval"))
		if err != nil {
			http.Error(w, "Invalid interval query parameter", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | query parameter are not valid | error: %s", err)
		}
		sampling, err := mapQueryParamToDuration(query.Get("sampling"))
		if err != nil {
			http.Error(w, "Invalid sampling query parameter", http.StatusBadRequest)
			log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | sampling parameter are not valid | error: %s", err)
		}

		entries, err := service.ReturnTemperatureChartData(interval, sampling)
		respBody := mapTimeStampToSpecifiedFormatForTemp(entries)

		if err != nil {
			http.Error(w, "Internal Server Error – ReturnTemperatureChartData failed", http.StatusInternalServerError)
			log.Errorf("[ERROR] api.web.ReturnTemperatureChartData | ReturnTemperatureChartData failed | err: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)

		log.Info("[END] api.web.ReturnTemperatureChartData")
	}
}

func mapQueryParamToF64(floatS string) (float64, error) {
	return strconv.ParseFloat(floatS, 64)
}
func mapQueryParamToDuration(durationS string) (time.Duration, error) {
	return time.ParseDuration(durationS)
}
func mapTimeStampToSpecifiedFormatForTemp(entries []TemperatureChartDataEntry) []TemperatureChartDataEntry {
	for i := 0; i < len(entries); i++ {
		t, _ := time.Parse(time.RFC3339Nano, entries[i].Timestamp)
		entries[i].Timestamp = t.Format("15:04:05")
	}

	return entries
}
func mapTimeStampToSpecifiedFormatForHumidity(entries []HumidityChartDataEntry) []HumidityChartDataEntry {
	for i := 0; i < len(entries); i++ {
		t, _ := time.Parse(time.RFC3339Nano, entries[i].Timestamp)
		entries[i].Timestamp = t.Format("15:04:05")
	}

	return entries
}
func mapTimeStampToSpecifiedFormatForMoisture(entries []MoistureChartDataEntry) []MoistureChartDataEntry {
	for i := 0; i < len(entries); i++ {
		t, _ := time.Parse(time.RFC3339Nano, entries[i].Timestamp)
		entries[i].Timestamp = t.Format("15:04:05")
	}

	return entries
}
