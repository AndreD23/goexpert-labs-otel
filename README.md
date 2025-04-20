# API de Temperatura por CEP

Esta aplicação é composta por dois serviços que trabalham em conjunto para fornecer a temperatura atual de uma cidade com base no CEP fornecido:

1. **Serviço A (Validação)** - Porta 8081
    - Responsável por receber e validar o CEP
    - Encaminha a requisição para o Serviço B

2. **Serviço B (Temperatura)** - Porta 8080
    - Utiliza a API do ViaCEP para obter a cidade
    - Consulta a WeatherAPI para obter a temperatura

O sistema utiliza o Zipkin para rastreamento distribuído das requisições.


---

## 🚀 Como Executar

### Pré-requisitos
- Docker
- Docker Compose
- Chave de API da WeatherAPI (obtenha gratuitamente em [weatherapi.com](https://www.weatherapi.com))

### Configuração

1. Clone o repositório:
```bash
git clone https://github.com/AndreD23/goexpert-labs-otel.git
cd goexpert-labs-otel
```

2. Crie o arquivo de configuração `.env`:
```bash
cp serviceb/.env.example serviceb/.env
```

3. Edite o arquivo `.env` e adicione sua chave API:
```env
WEATHER_API_KEY=sua_chave_aqui
```

### Executando com Docker Compose

1. Construa e inicie os containers:
```bash
docker compose up --build
```

Os serviços estarão disponíveis em:
- Serviço A: `http://localhost:8081`
- Serviço B: `http://localhost:8080`
- Zipkin: `http://localhost:9411`


---
## 📌 Endpoints

### POST / (Serviço A - 8081)
Endpoint principal para obter a temperatura. Envie um POST com o CEP no corpo da requisição:
```bash
curl --location '[http://localhost:8081](http://localhost:8081)'
--header 'Content-Type: application/json'
--data '{ "zipcode": "05187010" }'
```

### GET /{cep} (Serviço B - 8080)
Endpoint interno utilizado pelo Serviço A para consultar a temperatura.
Retorna a temperatura atual da cidade correspondente ao CEP.

#### Exemplos de Respostas

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

## 🔍 Monitoramento com Zipkin

O Zipkin está disponível em `http://localhost:9411` e permite:
- Visualizar o trace completo das requisições
- Monitorar o tempo de resposta de cada serviço
- Identificar gargalos e falhas na comunicação entre os serviços

Para acessar os traces:
1. Abra `http://localhost:9411`
2. Clique em "Run Query" no menu superior
3. Ajuste os filtros conforme necessário
4. Clique em "SHOW"

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

## 🔧 Stack Tecnológica

- Go 1.23.8
- Docker
- Chi Router
- WeatherAPI
- ViaCEP
- Zipkin
- OpenTelemetry


---

## 📝 Notas

- O CEP deve conter 8 dígitos (apenas números)
- A API remove automaticamente caracteres especiais do CEP (como hífen)
- As temperaturas são retornadas em graus Celsius, Fahrenheit e Kelvin
- O sistema utiliza tracing distribuído para monitoramento de performance e debugging

