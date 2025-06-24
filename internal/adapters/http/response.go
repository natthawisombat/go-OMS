package adap_http

import (
	"fmt"
	"strconv"

	"go-oms/internal/domain/entities"
	pkglogger "go-oms/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, response entities.Response, options ...map[string]interface{}) error {
	ctx := c.UserContext()
	logger := pkglogger.FromContext(ctx)

	fields := []interface{}{
		"status_code", response.StatusCode,
		"status", response.Status,
	}

	if len(options) > 0 {
		for k, v := range options[0] {
			fields = append(fields, k, v)
		}
	}

	if response.Status == "OK" {
		logger.Infow("response success", fields...)
	} else {
		fields = append(fields, "error_message", response.ErrorMessage)
		logger.Errorw("response error", fields...)
	}

	response.TransactionCode = fmt.Sprintf("%s", ctx.Value(entities.RequestId))
	statusCode := response.StatusCode
	response.StatusCode = 0
	return c.Status(statusCode).JSON(response)
}

func ResponseCreateOrder(c *fiber.Ctx, order entities.Order) error {
	ctx := c.UserContext()
	logger := pkglogger.FromContext(ctx)

	fields := []interface{}{
		"status_code", 200,
		"status", "OK",
	}
	logger.Infow("response success", fields...)
	respCreateOrder := entities.ResponseOrder{
		OrderID:      order.ID,
		CustomerName: order.CustomerName,
		TotalAmount:  order.TotalAmount,
	}
	for _, v := range order.OrderItems {
		respCreateOrder.Items = append(respCreateOrder.Items, entities.ResponseItem{
			ProductName: v.ProductName,
			Quantity:    v.Quantity,
			Price:       v.Price,
		})
	}

	return c.Status(200).JSON(respCreateOrder)
}

func ResponseGetOrder(c *fiber.Ctx, orders []entities.Order, total int64) error {
	ctx := c.UserContext()
	logger := pkglogger.FromContext(ctx)

	fields := []interface{}{
		"status_code", 200,
		"status", "OK",
	}
	logger.Infow("response success", fields...)

	return c.Status(200).JSON(entities.ResponseGetOrder{
		Total:  int(total),
		Orders: orders,
	})
}

func ResponseUpdateOrder(c *fiber.Ctx, orderId string, status string) error {
	ctx := c.UserContext()
	logger := pkglogger.FromContext(ctx)

	fields := []interface{}{
		"status_code", 200,
		"status", "OK",
	}
	logger.Infow("response success", fields...)

	toint, _ := strconv.Atoi(orderId)
	return c.Status(200).JSON(entities.ResponseUpdateOrder{
		OrderId: toint,
		Status:  status,
	})
}
