package main

import (
	"go-frete/api/internal/domain"
	"go-frete/api/internal/handler"
	"go-frete/api/internal/infra"
	"go-frete/api/pkg/logger"
	"net/http"
)

func main() {
	log := logger.New()
	log.Info("Iniciando API de Conversão...")

	// 1. Inicia o MongoDB
	mongoURI := "mongodb://root:example@mongodb:27017"
	mongoAdapter, err := infra.NewMongoDBAdapter(mongoURI, "currency_db")
	if err != nil {
		log.Fatal("Falha ao conectar no MongoDB", "erro", err.Error())
	}
	log.Info("Conectado ao MongoDB com sucesso!")

	// 2. Inicia o Adapter da API Externa
	apiAdapter := infra.NewAwesomeAPIAdapter()

	// 3. Injeta TUDO no UseCase
	usecase := domain.NewConverterUseCase(apiAdapter, mongoAdapter, log)

	// 4. Injeta o UseCase no Handler
	httpHandler := handler.NewConverterHandler(usecase, log)

	http.HandleFunc("/converter", httpHandler.Handle)

	log.Info("Servidor rodando", "porta", 8080)
	http.ListenAndServe(":8080", nil)
}
