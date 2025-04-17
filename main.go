package main

import (
	"Solflora/api/esp"
	"Solflora/api/web"
	"Solflora/db"
	"Solflora/logger"
	"Solflora/state"
	sys "github.com/joho/godotenv"
	"net/http"
)

func main() {
	err := sys.Overload(".env.local", ".env.cloud")

	logger.Init()
	db.Init()

	var log = logger.Logger()
	if err != nil {
		log.Fatal("[ERROR] main() | failed to load .env.cloud / .env.local file | %s", err.Error())
	}

	// Database write testing
	//mock.Mock_db_population_from_state()

	modelState := state.NewModelState()
	deviceState := state.NewDeviceState()
	tuneState := state.NewTuneState()

	controlSamplingService := esp.NewControlSamplingService(modelState, deviceState, tuneState)
	controlHandlerService := web.NewControlHandlerService(deviceState, modelState, tuneState)

	http.HandleFunc("/api/esp", esp.ControlSampler(controlSamplingService))
	http.HandleFunc("/api/pump-water", web.WaterPumpControl(controlHandlerService))
	http.HandleFunc("/api/fan-control", web.AirFanControl(controlHandlerService))
	http.HandleFunc("/api/temp-control-sp", web.TemperatureSetPointControl(controlHandlerService))
	http.HandleFunc("/api/temp-sp", web.ReturnTemperatureSetPoint(controlHandlerService))
	http.HandleFunc("/api/temp-coef", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			web.ReturnTemperatureControlTuneProfile(controlHandlerService)(w, r)
		case http.MethodPost:
			web.SetTemperatureControlTuneProfile(controlHandlerService)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/temp-data", web.ReturnTemperatureChartData(controlHandlerService))
	http.HandleFunc("/api/humidity-data", web.ReturnHumidityChartData(controlHandlerService))
	http.HandleFunc("/api/moisture-data", web.ReturnMoistureChartData(controlHandlerService))

	log.Fatalf("[FATAL] main() | web server shut down | potential-err: %s\n",
		http.ListenAndServe(":8080", nil).Error())
}
