package domain

import (
	"errors"
	"go-frete/api/pkg/logger"
	"time"
)

// O contrato que a regra de negócio exige.
type RateProvider interface {
	GetRate(moeda string) (float64, error)
}

// A estrutura do Caso de Uso
type ConverterUseCase struct {
	provider RateProvider
	repo     ConversionRepository
	log      logger.Logger
}

type ConversionRecord struct {
	MoedaDestino    string    `bson:"currency" json:"currency"`
	Cotacao         float64   `bson:"cotacao" json:"cotacao"`
	ValorEntrada    float64   `bson:"valor_entrada" json:"valor_entrada"`
	ValorConvertido float64   `bson:"valor_convertido" json:"valor_convertido"`
	Data            time.Time `bson:"data" json:"data"`
}

type ConversionRepository interface {
	SaveHistory(record ConversionRecord) error
}

func NewConverterUseCase(p RateProvider, r ConversionRepository, l logger.Logger) *ConverterUseCase {
	return &ConverterUseCase{provider: p, repo: r, log: l}
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

	//Monta o registro de conversão para salvar no histórico
	record := ConversionRecord{
		MoedaDestino:    moeda,
		Cotacao:         cotacao,
		ValorEntrada:    valorBRL,
		ValorConvertido: valorConvertido,
		Data:            time.Now(),
	}

	if err := uc.repo.SaveHistory(record); err != nil {
		uc.log.Error("Falha ao salvar histórico no banco", "erro", err.Error())
		return 0, errors.New("erro interno ao salvar conversão")
	}

	uc.log.Info("Conversão finalizada com sucesso", "valor_convertido", valorConvertido)
	return valorConvertido, nil
}
