package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	WeatherAPIKey string
	WeatherAPIURL string
	ViaCEPURL     string
}

// LoadConfig loads the environment variables and returns a Config struct
func LoadConfig() (*Config, error) {
	// Load the .env file, but do not fail if it does not exist
	godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		WeatherAPIKey: getEnv("WEATHER_API_KEY", ""),
		WeatherAPIURL: getEnv("WEATHER_API_URL", "http://api.weatherapi.com/v1/current.json"),
		ViaCEPURL:     getEnv("VIA_CEP_URL", "http://viacep.com.br/ws"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
