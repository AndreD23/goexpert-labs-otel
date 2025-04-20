package handlers

import (
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
)

type RequestBody struct {
	ZipCode string `json:"zipcode"`
}

type ResponseServiceB struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type TemperatureHandler struct {
}

func New() *TemperatureHandler {
	return &TemperatureHandler{}
}

func (t *TemperatureHandler) validateZipCode(zipCode string) (string, error) {
	cleanZip := ""
	for _, char := range zipCode {
		if char >= '0' && char <= '9' {
			cleanZip += string(char)
		}
	}
	if len(cleanZip) != 8 {
		return "", fmt.Errorf("invalid zipcode: must contain exactly 8 digits")
	}
	return cleanZip, nil
}

func (t *TemperatureHandler) HandleZipCodeInput(w http.ResponseWriter, r *http.Request) {
	tr := otel.GetTracerProvider().Tracer("HandleZipCodeInput")

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := tr.Start(ctx, "zipcode validation")
	defer span.End()

	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cleanZip, err := t.validateZipCode(reqBody.ZipCode)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://appb:8080/"+cleanZip, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	respTemp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	respBody, err := io.ReadAll(respTemp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read temp body response: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
