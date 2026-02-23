FROM golang:alpine

WORKDIR /app

# Instalação do Air
RUN go install github.com/air-verse/air@latest

# Dependências do projeto
COPY go.mod go.sum ./
RUN go mod download

# Configuração do Air
COPY .air.toml ./

CMD ["air"]