package main

import (
	"log"

	"github.com/Victor-armando18/computations-engine-poc/internal/application/usecase"
	"github.com/Victor-armando18/computations-engine-poc/internal/infrastructure/engine"
	"github.com/Victor-armando18/computations-engine-poc/internal/infrastructure/repository"
	web "github.com/Victor-armando18/computations-engine-poc/internal/interfaces/http"
	"github.com/labstack/echo/v4"
)

func main() {
	// 1. Inicializa Infraestrutura (Borda)
	repo := repository.NewMemoryRepository() // Implementado anteriormente
	engAdapter, err := engine.NewAdapter("rules/order.rules.json")
	if err != nil {
		log.Fatal("Failed to load rules:", err)
	}

	// 2. Inicializa Domínio/Aplicação (Núcleo)
	orderUC := usecase.NewOrderUseCase(repo, engAdapter)

	// 3. Inicializa Interface de Usuário (HTTP)
	e := echo.New()
	handler := web.NewOrderHandler(orderUC)

	web.RegisterRoutes(e, handler)

	log.Println("Order Service started on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
