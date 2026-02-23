package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-frete/api/internal/domain"
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

func (m *repositoryMock) SaveHistory(record domain.ConversionRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

type conversionReaderMock struct {
	mock.Mock
}

func (m *conversionReaderMock) GetLastConversions(limit int) ([]domain.ConversionRecord, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ConversionRecord), args.Error(1)
}

type conversionSearcherMock struct {
	mock.Mock
}

func (m *conversionSearcherMock) GetConversionsByCurrency(currency string) ([]domain.ConversionRecord, error) {
	args := m.Called(currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ConversionRecord), args.Error(1)
}

func TestConverterHandler_Handle(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should return 200 OK with valid json",
			run:  shouldReturn200OkWithValidJson,
		},
		{
			name: "should return 400 Bad Request with invalid json body",
			run:  shouldReturn400BadRequestWithInvalidJson,
		},
		{
			name: "should return 405 Method Not Allowed for GET request",
			run:  shouldReturn405MethodNotAllowed,
		},
		{
			name: "should return 422 Unprocessable Entity when currency is not found",
			run:  shouldReturn422UnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldReturn200OkWithValidJson(t *testing.T) {
	providerMock := new(rateProviderMock)
	repoMock := new(repositoryMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "USD").Return(5.0, nil)
	repoMock.On("SaveHistory", mock.Anything).Return(nil)

	usecase := domain.NewConverterUseCase(providerMock, repoMock, loggerMock)
	listUseCase := domain.NewListConversionsUseCase(new(conversionReaderMock), loggerMock)
	handler := NewConverterHandler(usecase, listUseCase, nil, loggerMock)

	body := []byte(`{"moeda": "USD", "valor_brl": 100.0}`)
	req, _ := http.NewRequest(http.MethodPost, "/converter", bytes.NewBuffer(body))

	recorder := httptest.NewRecorder()
	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), `"valor_convertido":20`)
}

func shouldReturn400BadRequestWithInvalidJson(t *testing.T) {
	loggerMock := new(loggermock.LoggerMock)
	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Warn", mock.Anything, mock.Anything).Return()

	listUseCase := domain.NewListConversionsUseCase(new(conversionReaderMock), loggerMock)
	handler := NewConverterHandler(nil, listUseCase, nil, loggerMock)

	body := []byte(`{moeda: USD}`)
	req, _ := http.NewRequest(http.MethodPost, "/converter", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func shouldReturn405MethodNotAllowed(t *testing.T) {
	loggerMock := new(loggermock.LoggerMock)
	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Warn", mock.Anything, mock.Anything).Return()

	listUseCase := domain.NewListConversionsUseCase(new(conversionReaderMock), loggerMock)
	handler := NewConverterHandler(nil, listUseCase, nil, loggerMock)

	req, _ := http.NewRequest(http.MethodGet, "/converter", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func shouldReturn422UnprocessableEntity(t *testing.T) {
	providerMock := new(rateProviderMock)
	repoMock := new(repositoryMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Warn", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "XYZ").Return(0.0, errors.New("moeda_nao_encontrada"))

	usecase := domain.NewConverterUseCase(providerMock, repoMock, loggerMock)
	listUseCase := domain.NewListConversionsUseCase(new(conversionReaderMock), loggerMock)
	handler := NewConverterHandler(usecase, listUseCase, nil, loggerMock)

	body := []byte(`{"moeda": "XYZ", "valor_brl": 100.0}`)
	req, _ := http.NewRequest(http.MethodPost, "/converter", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestConverterHandler_ListHandle(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should return 200 OK with list of conversions",
			run:  shouldReturn200OkForList,
		},
		{
			name: "should return 405 Method Not Allowed for POST request",
			run:  shouldReturn405MethodNotAllowedForList,
		},
		{
			name: "should return 500 Internal Server Error when database fails",
			run:  shouldReturn500WhenDatabaseFailsForList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldReturn200OkForList(t *testing.T) {
	readerMock := new(conversionReaderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	readerMock.On("GetLastConversions", 10).Return(nil, nil)

	listUseCase := domain.NewListConversionsUseCase(readerMock, loggerMock)
	handler := NewConverterHandler(nil, listUseCase, nil, loggerMock)

	req, _ := http.NewRequest(http.MethodGet, "/convert/list", nil)
	recorder := httptest.NewRecorder()

	handler.ListHandle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func shouldReturn405MethodNotAllowedForList(t *testing.T) {
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Warn", mock.Anything, mock.Anything).Return()

	listUseCase := domain.NewListConversionsUseCase(new(conversionReaderMock), loggerMock)
	handler := NewConverterHandler(nil, listUseCase, nil, loggerMock)

	req, _ := http.NewRequest(http.MethodPost, "/convert/list", nil)
	recorder := httptest.NewRecorder()

	handler.ListHandle(recorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func shouldReturn500WhenDatabaseFailsForList(t *testing.T) {
	readerMock := new(conversionReaderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	expectedErr := errors.New("mongo timeout")
	readerMock.On("GetLastConversions", 10).Return(nil, expectedErr)

	listUseCase := domain.NewListConversionsUseCase(readerMock, loggerMock)
	handler := NewConverterHandler(nil, listUseCase, nil, loggerMock)

	req, _ := http.NewRequest(http.MethodGet, "/convert/list", nil)
	recorder := httptest.NewRecorder()

	handler.ListHandle(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestConverterHandler_VariationHandle(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "should return 200 OK with variation data",
			run:  shouldReturn200OkForVariation,
		},
		{
			name: "should return 400 Bad Request when currency is missing",
			run:  shouldReturn400BadRequestForMissingCurrency,
		},
		{
			name: "should return 500 when database fails",
			run:  shouldReturn500ForVariationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func shouldReturn200OkForVariation(t *testing.T) {
	searcherMock := new(conversionSearcherMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	searcherMock.On("GetConversionsByCurrency", "USD").Return(nil, nil)

	variationUseCase := domain.NewVariationUseCase(searcherMock, loggerMock)
	handler := NewConverterHandler(nil, nil, variationUseCase, loggerMock)

	req, _ := http.NewRequest(http.MethodGet, "/variation/USD", nil)
	req.SetPathValue("moeda", "USD")

	recorder := httptest.NewRecorder()
	handler.VariationHandle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func shouldReturn400BadRequestForMissingCurrency(t *testing.T) {
	loggerMock := new(loggermock.LoggerMock)
	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Warn", mock.Anything, mock.Anything).Return()

	handler := NewConverterHandler(nil, nil, nil, loggerMock)

	// Sem passar o parâmetro de rota "moeda"
	req, _ := http.NewRequest(http.MethodGet, "/variation/", nil)
	recorder := httptest.NewRecorder()

	handler.VariationHandle(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func shouldReturn500ForVariationError(t *testing.T) {
	searcherMock := new(conversionSearcherMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	searcherMock.On("GetConversionsByCurrency", "EUR").Return(nil, errors.New("db error"))

	variationUseCase := domain.NewVariationUseCase(searcherMock, loggerMock)
	handler := NewConverterHandler(nil, nil, variationUseCase, loggerMock)

	req, _ := http.NewRequest(http.MethodGet, "/variation/EUR", nil)
	req.SetPathValue("moeda", "EUR")

	recorder := httptest.NewRecorder()
	handler.VariationHandle(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
