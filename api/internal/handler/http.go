package handler

import (
	"encoding/json"
	"fmt"
	"go-frete/api/internal/domain"
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
	converterUseCase *domain.ConverterUseCase
	listUseCase      *domain.ListConversionsUseCase
	variationUseCase *domain.VariationUseCase
	log              logger.Logger
}

func NewConverterHandler(uc *domain.ConverterUseCase, luc *domain.ListConversionsUseCase, vuc *domain.VariationUseCase, l logger.Logger) *ConverterHandler {
	return &ConverterHandler{converterUseCase: uc, listUseCase: luc, variationUseCase: vuc, log: l}
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
	valorConvertido, err := h.converterUseCase.Execute(req.Moeda, req.ValorBRL)

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

func (h *ConverterHandler) ListHandle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			h.log.Error("Recuperado de pânico no list handler", "detalhe", rec)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		}
	}()

	h.log.Info("Recebendo requisição de listagem", "endpoint", r.URL.Path, "metodo", r.Method)

	if r.Method != http.MethodGet {
		h.log.Warn("Método HTTP não permitido para listagem", "metodo_recebido", r.Method)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Chama a regra de negócio
	records, err := h.listUseCase.Execute()
	if err != nil {
		h.log.Error("Falha ao processar listagem na regra de negócio", "erro", err.Error())
		http.Error(w, "Erro ao buscar histórico", http.StatusInternalServerError)
		return
	}

	h.log.Info("Listagem finalizada com sucesso")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)
}

func (h *ConverterHandler) VariationHandle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			h.log.Error("Recuperado de pânico no variation handler", "detalhe", rec)
			http.Error(w, "Erro interno no servidor", http.StatusInternalServerError)
		}
	}()

	// Captura a variável {moeda} da URL
	moeda := r.PathValue("moeda")

	if moeda == "" {
		h.log.Warn("Moeda não informada na rota")
		http.Error(w, "Moeda deve ser informada na rota (ex: /variation/JPY)", http.StatusBadRequest)
		return
	}

	h.log.Info("Recebendo requisição de variação", "moeda", moeda)

	variations, err := h.variationUseCase.Execute(moeda)
	if err != nil {
		h.log.Error("Falha ao calcular variação", "erro", err.Error())
		http.Error(w, "Erro ao buscar variações", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(variations)
}
