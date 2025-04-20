package weatherapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestWeatherAPI_GetTempByCity(t *testing.T) {
	apiKey := "test_api_key"
	api := NewWeatherAPI(apiKey)

	tests := []struct {
		name      string
		city      string
		mockResp  *Response
		mockErr   error
		wantErr   bool
		expectedC float64
		expectedF float64
		expectedK float64
	}{
		{
			name: "valid city",
			city: "London",
			mockResp: &Response{Temperature: struct {
				TempC float64 `json:"temp_c"`
				TempF float64 `json:"temp_f"`
				TempK float64 `json:"temp_k"`
			}{TempC: 20, TempF: 68}},
			mockErr:   nil,
			wantErr:   false,
			expectedC: 20,
			expectedF: 68,
			expectedK: 293,
		},
		{
			name: "city with special characters",
			city: "São Paulo",
			mockResp: &Response{Temperature: struct {
				TempC float64 `json:"temp_c"`
				TempF float64 `json:"temp_f"`
				TempK float64 `json:"temp_k"`
			}{TempC: 25, TempF: 77}},
			mockErr:   nil,
			wantErr:   false,
			expectedC: 25,
			expectedF: 77,
			expectedK: 298,
		},
		{
			name:      "API returns error",
			city:      "InvalidCity",
			mockResp:  nil,
			mockErr:   errors.New("API error"),
			wantErr:   true,
			expectedC: 0,
			expectedK: 0,
			expectedF: 0,
		},
	}

	// Guardar o cliente HTTP original
	originalClient := http.DefaultClient

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar o mock do transport
			mockTransport := &mockTransport{err: tt.mockErr}

			if tt.mockResp != nil {
				jsonData, _ := json.Marshal(tt.mockResp)
				mockTransport.response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
				}
			} else if tt.mockErr == nil {
				mockTransport.response = &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}
			}

			// Substituir o cliente HTTP padrão
			http.DefaultClient = &http.Client{
				Transport: mockTransport,
			}

			// Restaurar o cliente HTTP original após o teste
			defer func() {
				http.DefaultClient = originalClient
			}()

			resp, err := api.GetTempByCity(tt.city)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedC, resp.Temperature.TempC)
				assert.Equal(t, tt.expectedF, resp.Temperature.TempF)
				assert.Equal(t, tt.expectedK, resp.Temperature.TempK)
			}
		})
	}
}
