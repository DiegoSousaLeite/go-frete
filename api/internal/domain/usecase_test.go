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

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) SaveHistory(record ConversionRecord) error {
	args := m.Called(record)
	return args.Error(0)
}
func TestConverterUseCase_Execute(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should calculate conversion and save successfully",
			run:  shouldCalculateConversionAndSaveSuccessfully,
		},
		{
			name: "should return error when provider fails",
			run:  shouldReturnErrorWhenProviderFails,
		},
		{
			name: "should return error when rate is zero",
			run:  shouldReturnErrorWhenRateIsZero,
		},
		{
			name: "should return error when repository fails to save",
			run:  shouldReturnErrorWhenRepositoryFailsToSave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldCalculateConversionAndSaveSuccessfully(t *testing.T) {
	providerMock := new(rateProviderMock)
	repoMock := new(repositoryMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	providerMock.On("GetRate", "USD").Return(5.0, nil)

	repoMock.On("SaveHistory", mock.Anything).Return(nil)

	uc := NewConverterUseCase(providerMock, repoMock, loggerMock)
	result, err := uc.Execute("USD", 100.0)

	assert.NoError(t, err)
	assert.Equal(t, 20.0, result)

	providerMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
}

func shouldReturnErrorWhenProviderFails(t *testing.T) {
	providerMock := new(rateProviderMock)
	repoMock := new(repositoryMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	expectedErr := errors.New("api_error")
	providerMock.On("GetRate", "EUR").Return(0.0, expectedErr)

	uc := NewConverterUseCase(providerMock, repoMock, loggerMock)
	result, err := uc.Execute("EUR", 100.0)

	assert.Error(t, err)
	assert.Equal(t, 0.0, result)
	assert.Equal(t, expectedErr, err)

	providerMock.AssertExpectations(t)
	repoMock.AssertNotCalled(t, "SaveHistory")
}

func shouldReturnErrorWhenRateIsZero(t *testing.T) {
	providerMock := new(rateProviderMock)
	repoMock := new(repositoryMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "BTC").Return(0.0, nil)

	uc := NewConverterUseCase(providerMock, repoMock, loggerMock)
	_, err := uc.Execute("BTC", 100.0)

	assert.Error(t, err)
	assert.Equal(t, "cotação não pode ser zero", err.Error())

	providerMock.AssertExpectations(t)
	repoMock.AssertNotCalled(t, "SaveHistory")
}

func shouldReturnErrorWhenRepositoryFailsToSave(t *testing.T) {
	providerMock := new(rateProviderMock)
	repoMock := new(repositoryMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return() // Logará o erro do banco

	providerMock.On("GetRate", "USD").Return(5.0, nil)

	repoMock.On("SaveHistory", mock.Anything).Return(errors.New("mongo timeout"))

	uc := NewConverterUseCase(providerMock, repoMock, loggerMock)
	result, err := uc.Execute("USD", 100.0)

	assert.Error(t, err)
	assert.Equal(t, 0.0, result)
	assert.Equal(t, "erro interno ao salvar conversão", err.Error())

	providerMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
}
