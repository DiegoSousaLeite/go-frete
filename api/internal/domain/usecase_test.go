package domain

import (
	"errors"
	"testing"

	"go-frete/api/tests/mocks/loggermock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type rateProviderMock struct {
	mock.Mock
}

func (m *rateProviderMock) GetRate(moeda string) (float64, error) {
	args := m.Called(moeda)
	return args.Get(0).(float64), args.Error(1)
}

func TestConverterUseCase_Execute(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should calculate conversion successfully",
			run:  shouldCalculateConversionSuccessfully,
		},
		{
			name: "should return error when provider fails",
			run:  shouldReturnErrorWhenProviderFails,
		},
		{
			name: "should return error when rate is zero",
			run:  shouldReturnErrorWhenRateIsZero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldCalculateConversionSuccessfully(t *testing.T) {
	providerMock := new(rateProviderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "USD").Return(5.0, nil)

	uc := NewConverterUseCase(providerMock, loggerMock)
	result, err := uc.Execute("USD", 100.0)

	assert.NoError(t, err)
	assert.Equal(t, 20.0, result)
	providerMock.AssertExpectations(t)
}

func shouldReturnErrorWhenProviderFails(t *testing.T) {
	providerMock := new(rateProviderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	expectedErr := errors.New("api_error")
	providerMock.On("GetRate", "EUR").Return(0.0, expectedErr)

	uc := NewConverterUseCase(providerMock, loggerMock)
	result, err := uc.Execute("EUR", 100.0)

	assert.Error(t, err)
	assert.Equal(t, 0.0, result)
	assert.Equal(t, expectedErr, err)
	providerMock.AssertExpectations(t)
}

func shouldReturnErrorWhenRateIsZero(t *testing.T) {
	providerMock := new(rateProviderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "BTC").Return(0.0, nil)

	uc := NewConverterUseCase(providerMock, loggerMock)
	_, err := uc.Execute("BTC", 100.0)

	assert.Error(t, err)
	assert.Equal(t, "cotação não pode ser zero", err.Error())
	providerMock.AssertExpectations(t)
}
