package viacep

import (
	"github.com/AndreD23/goexpert-labs-otel/serviceb/pkg/utils"
)

type CepData struct {
	Localidade string `json:"localidade"`
}

type ViaCEP struct {
	CepData
}

func GetCityByZipCode(zipCode string) (string, error) {
	url := "https://viacep.com.br/ws/" + zipCode + "/json/"
	var data ViaCEP
	err := utils.FetchData(url, &data)
	if err != nil {
		return "", err
	}

	return data.Localidade, nil
}
