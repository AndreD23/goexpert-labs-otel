package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AndreD23/goexpert-labs-otel/serviceb/internal/weatherapi"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateZipCode(t *testing.T) {
	tests := []struct {
		name      string
		zipCode   string
		want      string
		expectErr bool
	}{
		{
			name:      "valid 8 digits",
			zipCode:   "12345678",
			want:      "12345678",
			expectErr: false,
		},
		{
			name:      "contains letters",
			zipCode:   "1234abcd",
			want:      "",
			expectErr: true,
		},
		{
			name:      "too short",
			zipCode:   "12345",
			want:      "",
			expectErr: true,
		},
		{
			name:      "too long",
			zipCode:   "123456789",
			want:      "",
			expectErr: true,
		},
		{
			name:      "contains special chars",
			zipCode:   "1234#678",
			want:      "",
			expectErr: true,
		},
		{
			name:      "contains dash char",
			zipCode:   "12345-678",
			want:      "12345678",
			expectErr: false,
		},
		{
			name:      "empty string",
			zipCode:   "",
			want:      "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := &TemperatureHandler{}
			got, err := th.validateZipCode(tt.zipCode)

			if (err != nil) != tt.expectErr {
				t.Errorf("validateZipCode() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateZipCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Mock do ViaCEP Service
type mockViaCEPService struct {
	mockResponse string
	mockError    error
}

// Mock do WeatherAPI
type mockWeatherAPI struct {
	mockResponse weatherapi.Response
	mockError    error
}

func (m *mockViaCEPService) GetCityByZipCode(zipCode string) (string, error) {
	return m.mockResponse, m.mockError
}

func (m *mockWeatherAPI) GetTempByCity(city string) (weatherapi.Response, error) {
	return m.mockResponse, m.mockError
}

func setupRouter(handler *TemperatureHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/temperature/{zipCode}", handler.GetTemperature)
	return r
}

func TestGetTemperature(t *testing.T) {
	tests := []struct {
		name             string
		zipCode          string
		mockCityResponse string
		mockCityError    error
		mockWeatherResp  weatherapi.Response
		mockWeatherError error
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Success",
			zipCode:          "12345678",
			mockCityResponse: "São Paulo",
			mockCityError:    nil,
			mockWeatherResp: weatherapi.Response{
				Temperature: struct {
					TempC float64 `json:"temp_c"`
					TempF float64 `json:"temp_f"`
					TempK float64 `json:"temp_k"`
				}{
					TempC: 25.0,
					TempF: 77.0,
					TempK: 298.0,
				},
			},
			mockWeatherError: nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"temp_c":25,"temp_f":77,"temp_k":298}`,
		},
		{
			name:             "Invalid ZipCode",
			zipCode:          "1234",
			mockCityResponse: "",
			mockCityError:    nil,
			mockWeatherResp:  weatherapi.Response{},
			mockWeatherError: nil,
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: "invalid zipcode",
		},
		{
			name:             "ViaCEP Error",
			zipCode:          "12345678",
			mockCityResponse: "",
			mockCityError:    errors.New("viacep error"),
			mockWeatherResp:  weatherapi.Response{},
			mockWeatherError: nil,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"viacep error"}`,
		},
		{
			name:             "Weather API Error",
			zipCode:          "12345678",
			mockCityResponse: "São Paulo",
			mockCityError:    nil,
			mockWeatherResp:  weatherapi.Response{},
			mockWeatherError: errors.New("weather api error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"weather api error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar mock do ViaCEP Service
			mockViaCEP := &mockViaCEPService{
				mockResponse: tt.mockCityResponse,
				mockError:    tt.mockCityError,
			}

			// Criar mock do WeatherAPI
			mockWeather := &mockWeatherAPI{
				mockResponse: tt.mockWeatherResp,
				mockError:    tt.mockWeatherError,
			}

			// Criar o handler com ambos os mocks
			handler := &TemperatureHandler{
				weatherAPI: mockWeather,
				viaCEP:     mockViaCEP,
			}

			// Setup router com Chi
			router := setupRouter(handler)

			// Criar request
			req := httptest.NewRequest("GET", "/temperature/"+tt.zipCode, nil)
			w := httptest.NewRecorder()

			// Executar request
			router.ServeHTTP(w, req)

			// Verificar status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar resposta
			if tt.expectedStatus == http.StatusOK {
				var got, want weatherapi.Response
				err := json.Unmarshal(w.Body.Bytes(), &got.Temperature)
				assert.NoError(t, err)
				err = json.Unmarshal([]byte(tt.expectedResponse), &want.Temperature)
				assert.NoError(t, err)
				assert.Equal(t, want.Temperature, got.Temperature)
			} else {
				assert.Contains(t, w.Body.String(), tt.expectedResponse)
			}
		})
	}
}
