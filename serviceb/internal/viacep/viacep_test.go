package viacep

import (
	"bytes"
	"encoding/json"
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

func TestGetCityByZipCode(t *testing.T) {
	// Guardar o cliente HTTP original
	originalClient := http.DefaultClient
	defer func() {
		http.DefaultClient = originalClient
	}()

	tests := []struct {
		name       string
		zipCode    string
		mockResp   *ViaCEP
		statusCode int
		expectErr  bool
		expected   string
	}{
		{
			name:    "ValidZipCode",
			zipCode: "12345678",
			mockResp: &ViaCEP{
				CepData: CepData{
					Localidade: "São Paulo",
				},
			},
			statusCode: http.StatusOK,
			expectErr:  false,
			expected:   "São Paulo",
		},
		{
			name:       "InvalidZipCode",
			zipCode:    "00000000",
			mockResp:   nil,
			statusCode: http.StatusNotFound,
			expectErr:  true,
			expected:   "",
		},
		{
			name:       "NetworkError",
			zipCode:    "99999999",
			mockResp:   nil,
			statusCode: http.StatusInternalServerError,
			expectErr:  true,
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := &mockTransport{}

			if tt.mockResp != nil {
				jsonData, _ := json.Marshal(tt.mockResp)
				mockTransport.response = &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
				}
			} else {
				mockTransport.response = &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewBufferString("")),
				}
			}

			http.DefaultClient = &http.Client{
				Transport: mockTransport,
			}

			city, err := GetCityByZipCode(tt.zipCode)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, city)
			}
		})
	}
}
