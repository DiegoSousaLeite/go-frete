package infra

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type AwesomeAPIData struct {
	Bid string `json:"bid"`
}

// O Adapter que implementa a Interface do Domain
type AwesomeAPIAdapter struct{}

func NewAwesomeAPIAdapter() *AwesomeAPIAdapter {
	return &AwesomeAPIAdapter{}
}

// GetRate cumpre o contrato exigido pelo domain.RateProvider
func (a *AwesomeAPIAdapter) GetRate(moeda string) (float64, error) {
	url := "https://economia.awesomeapi.com.br/json/last/" + moeda + "-BRL"

	resp, err := http.Get(url)
	if err != nil {
		return 0, errors.New("erro ao consultar cotação externa")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("erro ao ler resposta da API")
	}

	var apiResponse map[string]AwesomeAPIData
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return 0, errors.New("erro ao processar cotação")
	}

	mapKey := moeda + "BRL"
	data, ok := apiResponse[mapKey]
	if !ok {
		return 0, errors.New("moeda_nao_encontrada")
	}

	cotacao, err := strconv.ParseFloat(data.Bid, 64)
	if err != nil {
		return 0, errors.New("erro no valor da cotação")
	}

	return cotacao, nil
}
