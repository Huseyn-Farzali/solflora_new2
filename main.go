package main

import (
	"Solflora/api/esp"
	"Solflora/api/web"
	"Solflora/db"
	"Solflora/logger"
	"Solflora/state"
	"Solflora/util"
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
	integralState := state.NewTrackingIntegralState()

	controlSamplingService := esp.NewControlSamplingService(modelState, deviceState, tuneState, integralState)
	controlHandlerService := web.NewControlHandlerService(deviceState, modelState, tuneState)

	http.HandleFunc("/api/esp", util.WithCors(esp.ControlSampler(controlSamplingService)))
	http.HandleFunc("/api/pump-water", util.WithCors(web.WaterPumpControl(controlHandlerService)))
	http.HandleFunc("/api/fan-control", util.WithCors(web.AirFanControl(controlHandlerService)))
	http.HandleFunc("/api/temp-control-sp", util.WithCors(web.TemperatureSetPointControl(controlHandlerService)))
	http.HandleFunc("/api/temp-sp", util.WithCors(web.ReturnTemperatureSetPoint(controlHandlerService)))
	http.HandleFunc("/api/temp-coef", util.WithCors(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			web.ReturnTemperatureControlTuneProfile(controlHandlerService)(w, r)
		case http.MethodPost:
			web.SetTemperatureControlTuneProfile(controlHandlerService)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/api/temp-data", util.WithCors(web.ReturnTemperatureChartData(controlHandlerService)))
	http.HandleFunc("/api/humidity-data", util.WithCors(web.ReturnHumidityChartData(controlHandlerService)))
	http.HandleFunc("/api/moisture-data", util.WithCors(web.ReturnMoistureChartData(controlHandlerService)))

	log.Fatalf("[FATAL] main() | web server shut down | potential-err: %s\n",
		http.ListenAndServe(":8080", nil).Error())
}
