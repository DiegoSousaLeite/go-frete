package domain

import (
	"errors"
	"testing"
	"time"

	"go-frete/api/tests/mocks/loggermock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type conversionSearcherMock struct {
	mock.Mock
}

func (m *conversionSearcherMock) GetConversionsByCurrency(currency string) ([]ConversionRecord, error) {
	args := m.Called(currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ConversionRecord), args.Error(1)
}

func TestVariationUseCase_Execute(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should calculate variation successfully",
			run:  shouldCalculateVariationSuccessfully,
		},
		{
			name: "should return error when repository fails",
			run:  shouldReturnErrorWhenSearchFails,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldCalculateVariationSuccessfully(t *testing.T) {
	searcherMock := new(conversionSearcherMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	mockData := []ConversionRecord{
		{MoedaDestino: "JPY", Cotacao: 30.0, ValorEntrada: 100, ValorConvertido: 3000, Data: time.Now().Add(-1 * time.Hour)},
		{MoedaDestino: "JPY", Cotacao: 33.0, ValorEntrada: 100, ValorConvertido: 3300, Data: time.Now()}, // Aumentou 3.0 (10%)
	}

	searcherMock.On("GetConversionsByCurrency", "JPY").Return(mockData, nil)

	uc := NewVariationUseCase(searcherMock, loggerMock)
	result, err := uc.Execute("JPY")

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// O primeiro registo não tem variação (é o ponto de partida)
	assert.Equal(t, 0.0, result[0].VariacaoValor)
	assert.Equal(t, 0.0, result[0].VariacaoPercentual)

	// O segundo registo deve ter calculado a subida
	assert.Equal(t, 3.0, result[1].VariacaoValor)       // 33 - 30 = 3
	assert.Equal(t, 10.0, result[1].VariacaoPercentual) // (3 / 30) * 100 = 10%

	searcherMock.AssertExpectations(t)
}

func shouldReturnErrorWhenSearchFails(t *testing.T) {
	searcherMock := new(conversionSearcherMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	expectedErr := errors.New("db connection lost")
	searcherMock.On("GetConversionsByCurrency", "USD").Return(nil, expectedErr)

	uc := NewVariationUseCase(searcherMock, loggerMock)
	result, err := uc.Execute("USD")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	searcherMock.AssertExpectations(t)
}
