package main

import (
	"go-frete/api/domain"
	"go-frete/api/handler"
	"go-frete/api/infra"
	"net/http"
)

func main() {
	// 1. Cria a Infraestrutura
	apiAdapter := infra.NewAwesomeAPIAdapter()

	// 2. Cria o Caso de Uso, injetando o Adapter da API
	usecase := domain.NewConverterUseCase(apiAdapter)

	// 3. Cria o Handler HTTP, injetando o Caso de Uso nele
	httpHandler := handler.NewConverterHandler(usecase)

	// 4. Configura as rotas e sobe o servidor
	http.HandleFunc("/converter", httpHandler.Handle)

	println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
