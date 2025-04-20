package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchData(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		mockResponse  *http.Response
		mockError     error
		expectedError error
		expectedData  interface{}
	}{
		{
			name: "successful fetch",
			url:  "http://example.com/data",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"key":"value"}`))),
			},
			expectedError: nil,
			expectedData:  map[string]interface{}{"key": "value"},
		},
		{
			name: "invalid JSON response",
			url:  "http://example.com/invalid",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`invalid-json`))),
			},
			expectedError: &json.SyntaxError{},
			expectedData:  nil,
		},
		{
			name: "HTTP error received",
			url:  "http://example.com/404",
			mockResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewReader([]byte(``))),
			},
			expectedError: errors.New("404 Not Found"),
			expectedData:  nil,
		},
		{
			name:          "network error",
			url:           "http://example.com/error",
			mockResponse:  nil,
			mockError:     errors.New("network error"),
			expectedError: errors.New("network error"),
			expectedData:  nil,
		},
	}

	originalHttpClient := http.DefaultClient
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &http.Client{
				Transport: &mockTransport{
					mockResponse: tt.mockResponse,
					mockError:    tt.mockError,
				},
			}
			http.DefaultClient = mockClient
			defer func() { http.DefaultClient = originalHttpClient }()

			var data map[string]interface{}
			err := FetchData(tt.url, &data)

			if tt.expectedError != nil {
				assert.Error(t, err)
				var syntaxError *json.SyntaxError
				if errors.As(tt.expectedError, &syntaxError) {
					assert.IsType(t, &json.SyntaxError{}, err)
				} else if tt.name == "network error" {
					assert.Contains(t, err.Error(), "network error")
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}

type mockTransport struct {
	mockResponse *http.Response
	mockError    error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}
