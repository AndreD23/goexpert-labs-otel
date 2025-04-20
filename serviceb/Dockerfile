FROM golang:1.23 as build

# Criar e definir o diretório de trabalho
WORKDIR /go/src/app

# Copiar os arquivos go.mod e go.sum primeiro
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o resto do código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/temperature ./cmd/...

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=build /go/bin/temperature /app/
ENTRYPOINT ["/app/temperature"]
