# Go Training Challenges

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

