package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/xavierpms/weather-by-city/internal/config"
	"github.com/xavierpms/weather-by-city/internal/infra/repository"
	"github.com/xavierpms/weather-by-city/internal/infra/validator"
	"github.com/xavierpms/weather-by-city/internal/infra/webserver/handlers"
	"github.com/xavierpms/weather-by-city/internal/usecase"
)

func main() {
	// Load the configurations from the .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize the router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Inject dependencies
	cepValidator := validator.NewCEPValidator()
	cepRepository := repository.NewCEPRepository(cfg.ViaCEPURL)
	tempRepository := repository.NewTemperatureRepository(cfg.WeatherAPIURL, cfg.WeatherAPIKey)
	getTempUseCase := usecase.NewGetTemperatureByCEP(cepRepository, tempRepository, cepValidator)
	temperatureHandler := handlers.NewTemperatureHandler(getTempUseCase)

	// Define the routes
	router.Get("/{cep}", temperatureHandler.GetTemperatureByCEP)

	// Start the server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
