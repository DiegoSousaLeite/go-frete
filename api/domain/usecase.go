package domain

import (
	"errors"
	"go-frete/api/pkg/logger"
)

// O contrato que a regra de negócio exige.
type RateProvider interface {
	GetRate(moeda string) (float64, error)
}

// A estrutura do Caso de Uso
type ConverterUseCase struct {
	provider RateProvider
	log      logger.Logger
}

func NewConverterUseCase(p RateProvider, l logger.Logger) *ConverterUseCase {
	return &ConverterUseCase{provider: p, log: l}
}

// A Regra de Negócio Pura
func (uc *ConverterUseCase) Execute(moeda string, valorBRL float64) (float64, error) {

	uc.log.Info("Iniciando cálculo de conversão",
		"moeda_alvo", moeda,
		"valor_brl", valorBRL,
	)
	// 1. Pede a cotação pro "provedor"
	cotacao, err := uc.provider.GetRate(moeda)
	if err != nil {
		uc.log.Error("Falha ao buscar cotação no provider", "erro", err.Error())
		return 0, err
	}

	if cotacao == 0 {
		return 0, errors.New("cotação não pode ser zero")
	}

	// 2. Faz a matemática
	valorConvertido := valorBRL / cotacao

	uc.log.Info("Conversão finalizada com sucesso", "valor_convertido", valorConvertido)
	return valorConvertido, nil
}
