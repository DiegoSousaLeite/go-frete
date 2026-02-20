package main

import (
	"go-frete/api/domain"
	"go-frete/api/handler"
	"go-frete/api/infra"
	"go-frete/api/pkg/logger"
	"net/http"
)

func main() {

	// 1. Inicia o Logger Global
	log := logger.New()
	log.Info("Iniciando API de Conversão...")

	// 2. Injeta as dependências
	apiAdapter := infra.NewAwesomeAPIAdapter()

	usecase := domain.NewConverterUseCase(apiAdapter, log)

	httpHandler := handler.NewConverterHandler(usecase, log)
	http.HandleFunc("/converter", httpHandler.Handle)

	log.Info("Servidor rodando", "porta", 8080)
	http.ListenAndServe(":8080", nil)
}
