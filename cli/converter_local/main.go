package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Campos com letra Maiúscula são obrigatórios para o JSON conseguir ler.
type Conversion struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func main() {
	// 1. Validação de Argumentos
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: currency-converter <amount> <currency>")
		os.Exit(1)
	}

	// 2. Leitura de Arquivo (I/O)
	file, err := os.ReadFile("../rates.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		os.Exit(1)
	}

	// 3. Parse do JSON (Unmarshal)
	var conversion Conversion

	err = json.Unmarshal(file, &conversion)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	// 4. Conversão de Tipos
	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid amount:", args[0])
		os.Exit(1)
	}

	// 5. Busca Segura no Mapa
	tax, ok := conversion.Rates[args[1]]

	if !ok {
		fmt.Fprintln(os.Stderr, "Currency not found:", args[1])
		os.Exit(1)
	}

	// 6. Cálculo e Saída
	result := amount * tax
	fmt.Printf("%.2f %s\n", result, args[1])
}
