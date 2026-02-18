package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// CurrencyData representa os dados de uma única moeda retornados pela API.
type CurrencyData struct {
	Bid string `json:"bid"`
}

func main() {
	// 1. Validação dos Argumentos
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Uso correto: ./convert <valor> <moeda>")
		fmt.Fprintln(os.Stderr, "Exemplo: ./convert 10 USD")
		os.Exit(1)
	}

	// 2. Parse do Valor (String -> Float)
	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Valor inválido:", args[0])
		os.Exit(1)
	}

	targetCurrency := args[1]

	// 3. Construção da URL e Requisição HTTP
	url := "https://economia.awesomeapi.com.br/json/last/" + targetCurrency + "-BRL"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao acessar API:", err)
		os.Exit(1)
	}
	// Fecha o corpo da resposta quando a função termina
	defer resp.Body.Close()

	// 4. Leitura do Corpo da Resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao ler resposta:", err)
		os.Exit(1)
	}

	// 5. Unmarshal do JSON Dinâmico
	var apiResponse map[string]CurrencyData

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao processar JSON:", err)
		os.Exit(1)
	}

	// 6. Extração dos Dados
	mapKey := targetCurrency + "BRL"
	data, ok := apiResponse[mapKey]

	if !ok {
		fmt.Printf("Moeda '%s' não encontrada ou não suportada pela API.\n", targetCurrency)
		os.Exit(1)
	}

	// 7. Conversão da Taxa (String -> Float)
	rate, err := strconv.ParseFloat(data.Bid, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erro ao converter a taxa recebida:", err)
		os.Exit(1)
	}

	// 8. Cálculo Final e Exibição
	result := amount * 1 / rate

	fmt.Printf("R$ %.2f equivale a %.2f %s\n", amount, result, targetCurrency)
	fmt.Printf("Cotação usada: %s (Fonte: AwesomeAPI)\n", data.Bid)
}
