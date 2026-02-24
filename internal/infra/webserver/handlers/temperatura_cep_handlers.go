package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/xavierpms/weather-by-city/internal/domain"
)

// TemperatureHandler handle the requests related to temperature
type TemperatureHandler struct {
	useCase domain.TemperatureUseCase
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewTemperatureHandler creates a new temperature handler
func NewTemperatureHandler(useCase domain.TemperatureUseCase) *TemperatureHandler {
	return &TemperatureHandler{
		useCase: useCase,
	}
}

// GetTemperatureByCEP handles the GET /{cep} request
func (h *TemperatureHandler) GetTemperatureByCEP(w http.ResponseWriter, r *http.Request) {
	// Extract the CEP from the URL
	cep := chi.URLParam(r, "cep")

	// Execute the use case
	temperature, err := h.useCase.GetTemperatureByCEP(cep)

	// Handle errors
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperature)
}

// handleError handles errors and returns the appropriate response
func (h *TemperatureHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case domain.ErrInvalidCEPFormat:
		log.Printf("invalid zipcode: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})

	case domain.ErrCEPNotFound:
		log.Printf("CEP not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "can not find zipcode"})

	case domain.ErrTemperatureNotFound:
		log.Printf("temperature not found: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "can not fetch temperature"})

	default:
		log.Printf("internal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
	}
}