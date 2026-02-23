package domain

import (
	"errors"
	"testing"
	"time"

	"go-frete/api/tests/mocks/loggermock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type conversionReaderMock struct {
	mock.Mock
}

func (m *conversionReaderMock) GetLastConversions(limit int) ([]ConversionRecord, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ConversionRecord), args.Error(1)
}

func TestListConversionsUseCase_Execute(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should return list of conversions successfully",
			run:  shouldReturnListOfConversionsSuccessfully,
		},
		{
			name: "should handle empty database gracefully",
			run:  shouldHandleEmptyDatabaseGracefully,
		},
		{
			name: "should return error when repository fails",
			run:  shouldReturnErrorWhenRepositoryFailsToRead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldReturnListOfConversionsSuccessfully(t *testing.T) {
	readerMock := new(conversionReaderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	mockData := []ConversionRecord{
		{MoedaDestino: "USD", Cotacao: 5.0, ValorEntrada: 100, ValorConvertido: 20, Data: time.Now()},
		{MoedaDestino: "EUR", Cotacao: 6.0, ValorEntrada: 120, ValorConvertido: 20, Data: time.Now()},
	}

	readerMock.On("GetLastConversions", 10).Return(mockData, nil)

	uc := NewListConversionsUseCase(readerMock, loggerMock)
	result, err := uc.Execute()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "USD", result[0].MoedaDestino)
	readerMock.AssertExpectations(t)
}

func shouldHandleEmptyDatabaseGracefully(t *testing.T) {
	readerMock := new(conversionReaderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	readerMock.On("GetLastConversions", 10).Return(nil, nil)

	uc := NewListConversionsUseCase(readerMock, loggerMock)
	result, err := uc.Execute()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	readerMock.AssertExpectations(t)
}

func shouldReturnErrorWhenRepositoryFailsToRead(t *testing.T) {
	readerMock := new(conversionReaderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	expectedErr := errors.New("mongo timeout")
	readerMock.On("GetLastConversions", 10).Return(nil, expectedErr)

	uc := NewListConversionsUseCase(readerMock, loggerMock)
	result, err := uc.Execute()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	readerMock.AssertExpectations(t)
}
