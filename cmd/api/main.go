// package main

// import (
// 	"github.com/Victor-armando18/computations-engine-poc/internal/http/handler"
// 	"github.com/labstack/echo/v4"
// )

// func main() {
// 	e := echo.New()
// 	h := handler.New()

// 	e.POST("/orders/:id/patch", h.Patch)
// 	e.POST("/orders/:id/validate", h.Validate)

// 	e.Logger.Fatal(e.Start(":8080"))
// }

package main

import (
	"log"

	"github.com/Victor-armando18/computations-engine-poc/internal/application/usecase"
	"github.com/Victor-armando18/computations-engine-poc/internal/infrastructure/engine"
	"github.com/Victor-armando18/computations-engine-poc/internal/infrastructure/repository"
	httpiface "github.com/Victor-armando18/computations-engine-poc/internal/interfaces/http"
	"github.com/labstack/echo/v4"
)

func main() {
	repo := repository.NewOrderMemoryRepository()
	eng, err := engine.NewAdapter("rules/order.rules.json")
	if err != nil {
		log.Fatal(err)
	}

	patch := usecase.NewPatchOrder(repo, eng)

	e := echo.New()
	httpiface.Register(e, &httpiface.OrderHandler{
		Patch:    patch,
		Validate: patch,
	})

	e.Logger.Fatal(e.Start(":8080"))
}
