package weatherapi

type WeatherAPIInterface interface {
	GetTempByCity(city string) (Response, error)
}
