package routers

import (
	"go-oms/configs"
	adap_http "go-oms/internal/adapters/http"
	adap_repository "go-oms/internal/adapters/repository"
	"go-oms/internal/domain/entities"
	usecases_orders "go-oms/internal/usecases/orders"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(cfg *configs.Setting) {
	// validate := validator.New()
	prefix := cfg.App.Group(configs.App.Prefix)

	Orders(cfg, prefix)

	prefix.Get("/healthcheck", func(c *fiber.Ctx) error {
		return adap_http.Response(c, entities.Response{Status: "OK", Message: "Healthy"}, map[string]interface{}{"function": "Healthcheck"})
	})

	prefix.Use(func(c *fiber.Ctx) error {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorCode: "ER404", ErrorMessage: "ไม่พบ Path", StatusCode: 404})
	})

}

func Orders(cfg *configs.Setting, prefix fiber.Router) {
	validate := validator.New()
	ordersGroup := prefix.Group("/orders")

	repository := adap_repository.NewRepository(cfg.PG.Store)
	usecases := usecases_orders.NewOrders(repository, validate)

	ordersGroup.Post("/", usecases.CreateOrders)
	ordersGroup.Get("", usecases.GetOrders)
	ordersGroup.Get("/:id", usecases.GetOrders)
	ordersGroup.Put("/:id", usecases.UpdateOrder)
}
