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

	mongoURI := "mongodb://root:example@mongodb:27017"
	mongoAdapter, err := infra.NewMongoDBAdapter(mongoURI, "currency_db")
	if err != nil {
		log.Fatal("Falha ao conectar no MongoDB", "erro", err.Error())
	}
	log.Info("Conectado ao MongoDB com sucesso!")

	apiAdapter := infra.NewAwesomeAPIAdapter()

	// 1. Injeta os 3 Casos de Uso!
	usecase := domain.NewConverterUseCase(apiAdapter, mongoAdapter, log)
	listUseCase := domain.NewListConversionsUseCase(mongoAdapter, log)
	variationUseCase := domain.NewVariationUseCase(mongoAdapter, log)

	// 2. Injeta no Handler
	httpHandler := handler.NewConverterHandler(usecase, listUseCase, variationUseCase, log)

	// 3. Rotas com suporte a variáveis de Path
	http.HandleFunc("POST /converter", httpHandler.Handle)
	http.HandleFunc("GET /convert/list", httpHandler.ListHandle)
	http.HandleFunc("GET /variation/{moeda}", httpHandler.VariationHandle)

	log.Info("Servidor rodando", "porta", 8080)
	http.ListenAndServe(":8080", nil)
}
