# API de Temperatura por CEP

Esta API retorna a temperatura atual de uma cidade com base no CEP fornecido. <br>
O serviço utiliza a API do ViaCEP para obter a cidade e a WeatherAPI para obter a temperatura.

---

## 🚀 Como Executar

### Pré-requisitos
- Docker
- Docker Compose
- Chave de API da WeatherAPI (obtenha gratuitamente em [weatherapi.com](https://www.weatherapi.com))

### Configuração

1. Clone o repositório:
```bash
git clone https://github.com/AndreD23/goexpert-labs-otel/service.git
cd goexpert-labs-cloudrun
```

2. Crie o arquivo de configuração `.env`:
```bash
cp .env.example .env
```

3. Edite o arquivo `.env` e adicione sua chave API:
```env
WEATHER_API_KEY=sua_chave_aqui
```

### Executando com Docker Compose

1. Construa e inicie o container:
```bash
docker compose up --build
```

A API estará disponível em `http://localhost:8080`

---
## 📌 Endpoints

### GET /{cep}
Retorna a temperatura atual da cidade correspondente ao CEP.

#### Exemplos de Requisições

1. CEP Válido:
```bash
curl http://localhost:8080/05187010
```
Resposta (200 OK):
```json
{"temp_c":20.2,"temp_f":68.4,"temp_k":293.2}
```

2. CEP Inválido (formato incorreto):
```bash
curl http://localhost:8080/123
```

Resposta (422 Unprocessable Entity):
```
invalid zipcode
```

3. CEP Não Encontrado:
```bash
curl http://localhost:8080/00000000
```

Resposta (404 Not Found):
```
can not find zipcode
```

---

## 🧪 Testes

Para executar os testes automatizados:

```bash
# Executar todos os testes
go test ./... -v

# Executar testes com cobertura
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🌎 Versão Online

A API está disponível online através do Google Cloud Run: <br>
[https://goexpert-labs-cloudrun-611389022433.southamerica-east1.run.app/](https://goexpert-labs-cloudrun-611389022433.southamerica-east1.run.app/)

### Exemplo de uso online:
```bash
curl https://goexpert-labs-cloudrun-611389022433.southamerica-east1.run.app/05187010
```

---

## 📝 Notas

- O CEP deve conter 8 dígitos (apenas números)
- A API remove automaticamente caracteres especiais do CEP (como hífen)
- As temperaturas são retornadas em graus Celsius, Fahrenheit e Kelvin

## 🔧 Stack Tecnológica

- Go 1.23
- Docker
- Chi Router
- WeatherAPI
- ViaCEP
