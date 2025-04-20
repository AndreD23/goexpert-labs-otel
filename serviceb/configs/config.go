package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
}

func NewConfig() *Config {
	return config
}

func init() {
	var err error
	config, err = loadConfig()
	if err != nil {
		panic(fmt.Sprintf("Erro ao carregar configurações: %v", err))
	}
}

func loadConfig() (*Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Habilita o carregamento de variáveis de ambiente
	viper.AutomaticEnv()

	// Tenta ler o arquivo .env, mas ignora se não existir
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Arquivo .env não encontrado. Usando variáveis de ambiente.")
	}

	// Define valores padrão
	viper.SetDefault("WEATHER_API_KEY", "")

	// Cria uma nova instância de Config
	config := &Config{}

	// Tenta carregar variáveis de ambiente diretamente
	weatherAPIKey := viper.GetString("WEATHER_API_KEY")
	if weatherAPIKey != "" {
		config.WeatherAPIKey = weatherAPIKey
	}

	// Validação das configurações obrigatórias
	if config.WeatherAPIKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY é obrigatória")
	}

	return config, nil

}
