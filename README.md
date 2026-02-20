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

Uma API REST que recebe um valor em BRL e converte para moeda estrangeira.
Este projeto aplica o padrão **Hexagonal (Ports and Adapters)**, separando completamente as Regras de Negócio (Domínio) da Infraestrutura (HTTP e APIs externas) através de **Injeção de Dependência**.

### 📂 Estrutura

```text
api/
├── domain/            # 🟡 Regra de negócio pura e Interfaces (Contratos)
├── infra/             # 🔵 Adapters (Integração externa com a AwesomeAPI)
├── handlers/          # 🔵 Delivery (Recebe e responde requisições HTTP)
└── main.go            # ⚙️ Ponto de entrada e Injeção de Dependências
```

### ⚡ Como Rodar

1. Entre no diretório da API:

```bash
cd api
```

2. Inicie o servidor (ele rodará na porta `:8080`):

```bash
go run .
```

### 🧪 Como Testar

**Via cURL (Terminal):**

```bash
curl -X POST http://localhost:8080/converter \
     -H "Content-Type: application/json" \
     -d '{"moeda": "USD", "valor_brl": 100}'
```

**Via HTTP Client (Postman/Insomnia):**

* **Método:** `POST`
* **URL:** `http://localhost:8080/converter`
* **Body (JSON):**
```json
{
  "moeda": "EUR",
  "valor_brl": 150.50
}
```



### 🛠 Status Codes Implementados

* `200 OK`: Conversão realizada com sucesso.
* `400 Bad Request`: Corpo da requisição ausente ou JSON mal formatado.
* `405 Method Not Allowed`: Tentativa de acesso com método diferente de POST.
* `422 Unprocessable Entity`: Cotação da moeda solicitada não foi encontrada.
* `500 Internal Server Error / 502 Bad Gateway`: Falha no servidor ou na API externa.

### 🛡️ Testes Automatizados

O projeto conta com uma suíte de testes unitários focada em garantir a confiabilidade da aplicação, cobrindo as regras de negócio (Domain) e a camada de entrega (Handlers).

**Stack de Testes:**
* **`testing` & `httptest`**: Pacotes nativos do Go para testes de mesa (Table-Driven) e simulação de requisições HTTP sem a necessidade de instanciar um servidor real.
* **`testify/assert`**: Para asserções limpas e sem repetição de código.
* **`testify/mock`**: Utilizado para criação de *Strict Mocks* (Mocks Estritos) globais e locais, isolando o comportamento de integrações externas e utilitários (como o Logger e a API de cotação).

**Como rodar os testes:**

Para executar toda a suíte de testes com detalhes dos cenários (verbose), utilize o comando na raiz da pasta `api`:

```bash
go test ./... -v