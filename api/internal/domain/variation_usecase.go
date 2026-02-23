package domain

import (
	"go-frete/api/pkg/logger"
	"time"
)

// O DTO de resposta
type CurrencyVariation struct {
	Data               time.Time `json:"data"`
	Cotacao            float64   `json:"cotacao"`
	VariacaoValor      float64   `json:"variacao_valor"`
	VariacaoPercentual float64   `json:"variacao_percentual"`
}

type ConversionSearcher interface {
	GetConversionsByCurrency(currency string) ([]ConversionRecord, error)
}

type VariationUseCase struct {
	repo ConversionSearcher
	log  logger.Logger
}

func NewVariationUseCase(r ConversionSearcher, l logger.Logger) *VariationUseCase {
	return &VariationUseCase{repo: r, log: l}
}

func (uc *VariationUseCase) Execute(moeda string) ([]CurrencyVariation, error) {
	uc.log.Info("Iniciando cálculo de variação", "moeda", moeda)

	// Busca os registros dessa moeda (ordenando da mais antiga para a mais nova)
	records, err := uc.repo.GetConversionsByCurrency(moeda)
	if err != nil {
		uc.log.Error("Falha ao buscar conversões por moeda", "erro", err.Error(), "moeda", moeda)
		return nil, err
	}

	var variations []CurrencyVariation

	// Regra de Negócio: Calcular a variação entre uma cotação e a anterior
	for i, record := range records {
		variacaoValor := 0.0
		variacaoPerc := 0.0

		// Se não for o primeiro registro, compara com o anterior
		if i > 0 {
			cotacaoAnterior := records[i-1].Cotacao
			variacaoValor = record.Cotacao - cotacaoAnterior
			variacaoPerc = (variacaoValor / cotacaoAnterior) * 100
		}

		variations = append(variations, CurrencyVariation{
			Data:               record.Data,
			Cotacao:            record.Cotacao,
			VariacaoValor:      variacaoValor,
			VariacaoPercentual: variacaoPerc,
		})
	}

	uc.log.Info("Cálculo de variação finalizado com sucesso", "total_registros", len(variations))
	return variations, nil
}
