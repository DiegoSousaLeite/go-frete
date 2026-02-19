package handler

import (
	"encoding/json"
	"fmt"
	"go-frete/api/domain"
	"net/http"
)

type Request struct {
	Moeda    string  `json:"moeda"`
	ValorBRL float64 `json:"valor_brl"`
}

type Response struct {
	ValorConvertido float64 `json:"valor_convertido"`
}

// O "Garçom" que atende o cliente
type ConverterHandler struct {
	usecase *domain.ConverterUseCase
}

func NewConverterHandler(uc *domain.ConverterUseCase) *ConverterHandler {
	return &ConverterHandler{usecase: uc}
}

func (h *ConverterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recuperado de pânico:", r)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// CHAMA A REGRA DE NEGÓCIO
	valorConvertido, err := h.usecase.Execute(req.Moeda, req.ValorBRL)

	if err != nil {
		// Tratamento de erros customizados
		if err.Error() == "moeda_nao_encontrada" {
			http.Error(w, "Moeda não encontrada ou inválida", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// DEVOLVE A RESPOSTA
	respo := Response{ValorConvertido: valorConvertido}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respo)
}
