package repository

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/xavierpms/weather-by-city/internal/domain"
)

// WeatherAPIResponse represents the response from the WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

// TemperatureRepository implements domain.TemperatureRepository
type TemperatureRepository struct {
	apiURL string
	apiKey string
}

// NewTemperatureRepository creates a new temperature repository
func NewTemperatureRepository(apiURL, apiKey string) domain.TemperatureRepository {
	return &TemperatureRepository{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

// GetTemperatureByCityName fetches the temperature for a given city
func (r *TemperatureRepository) GetTemperatureByCityName(cityName string) (*domain.Temperature, error) {
	// Encode the city name for URL
	encodedCityName := url.QueryEscape(cityName)

	// Build the URL with parameters
	requestURL := r.apiURL + "?q=" + encodedCityName + "&lang=pt&country=Brazil&key=" + r.apiKey

	// Make the request
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	var weatherResp WeatherAPIResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return nil, err
	}

	// Calculate the temperature in Kelvin
	kelvin := weatherResp.Current.TempC + 273.0

	return &domain.Temperature{
		Celsius:    weatherResp.Current.TempC,
		Fahrenheit: weatherResp.Current.TempF,
		Kelvin:     kelvin,
	}, nil
}
