package domain

import "errors"

// O contrato que a regra de negócio exige.
type RateProvider interface {
	GetRate(moeda string) (float64, error)
}

// A estrutura do Caso de Uso
type ConverterUseCase struct {
	provider RateProvider
}

func NewConverterUseCase(p RateProvider) *ConverterUseCase {
	return &ConverterUseCase{provider: p}
}

// A Regra de Negócio Pura
func (uc *ConverterUseCase) Execute(moeda string, valorBRL float64) (float64, error) {

	// 1. Pede a cotação pro "provedor"
	cotacao, err := uc.provider.GetRate(moeda)
	if err != nil {
		return 0, err
	}

	if cotacao == 0 {
		return 0, errors.New("cotação não pode ser zero")
	}

	// 2. Faz a matemática
	valorConvertido := valorBRL / cotacao

	return valorConvertido, nil
}
