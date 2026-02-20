package handler

import (
	"encoding/json"
	"fmt"
	"go-frete/api/domain"
	"go-frete/api/pkg/logger"
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
	log     logger.Logger
}

func NewConverterHandler(uc *domain.ConverterUseCase, l logger.Logger) *ConverterHandler {
	return &ConverterHandler{usecase: uc, log: l}
}

func (h *ConverterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recuperado de pânico:", r)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		}
	}()

	h.log.Info("Recebendo requisição", "endpoint", r.URL.Path, "metodo", r.Method)
	if r.Method != http.MethodPost {
		h.log.Warn("Método HTTP não permitido", "metodo_recebido", r.Method)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Falha ao fazer parse do JSON", "erro", err.Error())
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	h.log.Info("Dados validados com sucesso", "moeda", req.Moeda, "valor_brl", req.ValorBRL)

	// CHAMA A REGRA DE NEGÓCIO
	valorConvertido, err := h.usecase.Execute(req.Moeda, req.ValorBRL)

	if err != nil {
		// Tratamento de erros customizados
		if err.Error() == "moeda_nao_encontrada" {
			h.log.Warn("Moeda solicitada não é suportada", "moeda_solicitada", req.Moeda)
			http.Error(w, "Moeda não encontrada ou inválida", http.StatusUnprocessableEntity)
			return
		}
		h.log.Error("Falha ao processar conversão na regra de negócio", "erro", err.Error(), "moeda", req.Moeda)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	h.log.Info("Requisição finalizada com sucesso", "valor_convertido", valorConvertido)

	// DEVOLVE A RESPOSTA
	respo := Response{ValorConvertido: valorConvertido}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respo)
}
