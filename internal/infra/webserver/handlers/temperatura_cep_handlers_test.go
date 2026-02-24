package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xavierpms/weather-by-city/internal/domain"
)

// MockTemperatureUseCase is a mock of the TemperatureUseCase for testing
type MockTemperatureUseCase struct {
	getTemperatureByCEPFunc func(cep string) (*domain.Temperature, error)
}

func (m *MockTemperatureUseCase) GetTemperatureByCEP(cep string) (*domain.Temperature, error) {
	return m.getTemperatureByCEPFunc(cep)
}

// TestGetTemperatureByCEPSuccess tests the success of the request
func TestGetTemperatureByCEPSuccess(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(cep string) (*domain.Temperature, error) {
			return &domain.Temperature{
				Celsius:    28.5,
				Fahrenheit: 83.3,
				Kelvin:     301.65,
			}, nil
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/32450000", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var temperature domain.Temperature
	err := json.Unmarshal(w.Body.Bytes(), &temperature)
	assert.NoError(t, err)
	assert.Equal(t, 28.5, temperature.Celsius)
	assert.Equal(t, 83.3, temperature.Fahrenheit)
}

// TestGetTemperatureByCEPInvalidFormat tests the case when the CEP has an invalid format
func TestGetTemperatureByCEPInvalidFormat(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(cep string) (*domain.Temperature, error) {
			return nil, domain.ErrInvalidCEPFormat
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/3245000000", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", errResponse.Message)
}

// TestGetTemperatureByCEPNotFound tests the case when the CEP is not found
func TestGetTemperatureByCEPNotFound(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(cep string) (*domain.Temperature, error) {
			return nil, domain.ErrCEPNotFound
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/00000000", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", errResponse.Message)
}

// TestGetTemperatureByCEPTemperatureNotFound tests the case when the temperature is not found
func TestGetTemperatureByCEPTemperatureNotFound(t *testing.T) {
	// Arrange
	mockUseCase := &MockTemperatureUseCase{
		getTemperatureByCEPFunc: func(cep string) (*domain.Temperature, error) {
			return nil, domain.ErrTemperatureNotFound
		},
	}

	handler := NewTemperatureHandler(mockUseCase)
	req := httptest.NewRequest("GET", "/32450000", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetTemperatureByCEP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errResponse ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResponse)
	assert.NoError(t, err)
	assert.Equal(t, "can not fetch temperature", errResponse.Message)
}
