package main

import (
	"github.com/Victor-armando18/computations-engine-poc/internal/http/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	h := handler.New()

	e.POST("/orders/:id/patch", h.Patch)
	e.POST("/orders/:id/validate", h.Validate)

	e.Logger.Fatal(e.Start(":8080"))
}

// package main

// import (
// 	"context"
// 	"fmt"

// 	engine "github.com/dolphin-sistemas/computations-engine"
// 	"github.com/dolphin-sistemas/computations-engine/core"
// 	"github.com/dolphin-sistemas/computations-engine/loader"
// )

// func main() {
// 	// Carregar RulePack
// 	rulePack, err := loader.LoadRulePackFromFile("rules.json")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Criar estado
// 	state := core.State{
// 		TenantID: "tenant-1",
// 		Items: []core.Item{
// 			{
// 				ID:     "item-1",
// 				Amount: 11,
// 				Fields: map[string]interface{}{
// 					"basePrice": 100.0,
// 				},
// 			},
// 		},
// 		Fields: make(map[string]interface{}),
// 		Totals: core.Totals{},
// 	}

// 	// Executar motor
// 	_, _, reasons, violations, rulesVersion, err := engine.RunEngine(
// 		context.Background(),
// 		state,
// 		rulePack,
// 		core.ContextMeta{
// 			TenantID: "tenant-1",
// 			UserID:   "user-1",
// 			Locale:   "pt-BR",
// 		},
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Usar resultados
// 	fmt.Printf("Version: %s\n", rulesVersion)
// 	fmt.Printf("Reasons: %+v\n", reasons)
// 	fmt.Printf("Violations: %+v\n", violations)
// }
