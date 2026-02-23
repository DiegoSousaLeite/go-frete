package domain

import (
	"go-frete/api/pkg/logger"
)

// Limitar a busca
const SearchLimit = 10

type ConversionReader interface {
	GetLastConversions(limit int) ([]ConversionRecord, error)
}

type ListConversionsUseCase struct {
	repo ConversionReader
	log  logger.Logger
}

func NewListConversionsUseCase(r ConversionReader, l logger.Logger) *ListConversionsUseCase {
	return &ListConversionsUseCase{repo: r, log: l}
}

func (uc *ListConversionsUseCase) Execute() ([]ConversionRecord, error) {
	uc.log.Info("Iniciando busca do histórico de conversões")

	limit := SearchLimit

	records, err := uc.repo.GetLastConversions(limit)
	if err != nil {
		uc.log.Error("Falha ao buscar últimas conversões no banco de dados", "erro", err.Error())
		return nil, err
	}

	// Garante que não retorne nulo se o banco estiver vazio
	if records == nil {
		records = []ConversionRecord{}
	}

	uc.log.Info("Busca de histórico finalizada com sucesso", "quantidade_encontrada", len(records))
	return records, nil
}
