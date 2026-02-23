# Go Training Challenges 🚀

Repositório destinado aos exercícios e projetos práticos da trilha de aprendizado em **Golang**.

---

## 1️⃣ Desafio 1: Currency Converter CLI

Uma CLI (Command Line Interface) para conversão de moedas (BRL -> Moeda Estrangeira).
Este exercício explora conceitos fundamentais como **Leitura de Arquivos**, **JSON Parsing**, **Maps**, **Structs** e **HTTP Requests**.

### 📂 Estrutura

```text
cli/
├── converter_local/  # Versão Offline (Lê taxas do arquivo rates.json)
├── converter_api/    # Versão Online (Consulta a AwesomeAPI em tempo real)
└── rates.json        # Base de dados estática para testes locais
```

### ⚡ Como Rodar

> **Nota:** Para garantir a leitura correta dos arquivos, execute os comandos de dentro da pasta `cli`.

1. Entre no diretório do desafio:

```bash
cd cli
```

2. Execute a versão desejada:

**Versão Local (Arquivo Fixo):**

```bash
go run converter_local/main.go 150 EUR
```

**Versão API (Cotação Real):**

```bash
go run converter_api/main.go 50 USD
```

---

## 2️⃣ Desafio 2: Currency Converter API (Clean Architecture)

Uma API REST completa para conversão de moedas, armazenamento de histórico e cálculo de variações ao longo do tempo.

Este projeto aplica o padrão **Hexagonal (Ports and Adapters)** e **SOLID** (Segregação de Interfaces), separando completamente as Regras de Negócio (Domínio) da Infraestrutura (HTTP, MongoDB e APIs externas) através de **Injeção de Dependência**.

### 📂 Estrutura

```text
api/
├── internal/
│   ├── domain/     # Regras de negócio (UseCases) e Contratos (Interfaces segregadas)
│   ├── infra/      # Adapters (Integração com AwesomeAPI e MongoDB)
│   └── handler/    # Delivery (Controladores HTTP)
├── pkg/
│   └── logger/     # Utilitários compartilhados (Wrapper do Zap Logger)
├── tests/
│   └── mocks/      # Mocks globais compartilhados para testes
├── docker-compose.yaml
├── Dockerfile
└── main.go         # Ponto de entrada (Montador de Dependências)
```

### ⚡ Como Rodar (Docker)

Como o projeto agora depende de um banco de dados MongoDB, a melhor forma de executá-lo é via Docker Compose.

1. Entre no diretório da API:

```bash
cd api
```

2. Suba a aplicação e o banco de dados (o servidor rodará na porta `:8080` com *hot-reload* via Air):

```bash
docker compose up -d --build
```

*(Para ver os logs do sistema e do banco, utilize `docker compose logs -f app`)*

### 🧪 Endpoints e Como Testar

#### 1. Realizar Conversão (`POST /converter`)

Converte um valor em BRL para a moeda solicitada e salva o histórico no banco de dados.

```bash
curl -X POST http://localhost:8080/converter \
     -H "Content-Type: application/json" \
     -d '{"moeda": "USD", "valor_brl": 100}'
```

#### 2. Listar Histórico (`GET /convert/list`)

Retorna as últimas 10 conversões realizadas e salvas no banco de dados.

```bash
curl -X GET http://localhost:8080/convert/list
```

#### 3. Calcular Variação (`GET /variation/{moeda}`)

Busca todo o histórico de conversões de uma moeda específica e calcula a variação financeira e percentual entre cada operação no tempo.

```bash
curl -X GET http://localhost:8080/variation/USD
```

### 🛠 Status Codes Implementados

* `200 OK`: Operação realizada com sucesso.
* `400 Bad Request`: Corpo da requisição ausente, JSON mal formatado ou moeda não informada na rota.
* `405 Method Not Allowed`: Tentativa de acesso com método HTTP incorreto.
* `422 Unprocessable Entity`: Cotação da moeda solicitada não foi encontrada na API externa.
* `500 Internal Server Error / 502 Bad Gateway`: Falha interna no servidor, no banco de dados (MongoDB) ou na API externa.

### 🛡️ Testes Automatizados (100% Coverage)

O projeto conta com uma suíte de testes unitários focada em garantir a confiabilidade da aplicação, cobrindo as regras de negócio (Domain) e a camada de entrega (Handlers), com **100% de cobertura na camada de aplicação**.

**Stack de Testes:**

* **`testing` & `httptest`**: Pacotes nativos do Go para testes de mesa (Table-Driven) e simulação de requisições HTTP (incluindo variáveis de path do Go 1.22+).
* **`testify/assert`**: Para asserções limpas e legíveis.
* **`testify/mock`**: Utilizado para criação de *Strict Mocks* globais e locais, isolando o comportamento de integrações externas (MongoDB, APIs, Logs).

**Como rodar os testes localmente:**

1. Executar todos os testes com detalhes dos cenários (verbose):

```bash
go test ./... -v
```

2. Gerar relatório de cobertura de código (Coverage):

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

*(Isso abrirá uma página HTML no seu navegador mostrando exatamente quais linhas de código foram cobertas pelos testes).*
