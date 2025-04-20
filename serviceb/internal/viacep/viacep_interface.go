package viacep

type ViaCEPInterface interface {
	GetCityByZipCode(zipCode string) (string, error)
}

type DefaultViaCEPService struct{}

func NewViaCEPService() ViaCEPInterface {
	return &DefaultViaCEPService{}
}

func (s *DefaultViaCEPService) GetCityByZipCode(zipCode string) (string, error) {
	return GetCityByZipCode(zipCode)
}
