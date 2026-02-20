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
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "USD").Return(5.0, nil)

	usecase := domain.NewConverterUseCase(providerMock, loggerMock)
	handler := NewConverterHandler(usecase, loggerMock)

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

	handler := NewConverterHandler(nil, loggerMock)

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

	handler := NewConverterHandler(nil, loggerMock)

	req, _ := http.NewRequest(http.MethodGet, "/converter", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func shouldReturn422UnprocessableEntity(t *testing.T) {
	providerMock := new(rateProviderMock)
	loggerMock := new(loggermock.LoggerMock)

	loggerMock.On("Info", mock.Anything, mock.Anything).Return()
	loggerMock.On("Warn", mock.Anything, mock.Anything).Return()
	loggerMock.On("Error", mock.Anything, mock.Anything).Return()

	providerMock.On("GetRate", "XYZ").Return(0.0, errors.New("moeda_nao_encontrada"))

	usecase := domain.NewConverterUseCase(providerMock, loggerMock)
	handler := NewConverterHandler(usecase, loggerMock)

	body := []byte(`{"moeda": "XYZ", "valor_brl": 100.0}`)
	req, _ := http.NewRequest(http.MethodPost, "/converter", bytes.NewBuffer(body))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}
