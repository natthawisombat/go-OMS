package usecases_orders

import (
	adap_http "go-oms/internal/adapters/http"
	"go-oms/internal/domain/entities"
	domainrepo "go-oms/internal/domain/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Orders struct {
	repo      domainrepo.Repository
	validator *validator.Validate
}

func NewOrders(repo domainrepo.Repository, validator *validator.Validate) *Orders {
	return &Orders{repo: repo, validator: validator}
}

func (uc *Orders) CreateOrders(c *fiber.Ctx) error {
	var order entities.Order
	if err := c.BodyParser(&order); err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "CreateOrders"})
	}

	// validate request body
	err := uc.validator.Struct(&order)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "CreateOrders"})
		}
	}

	if err := calcurateTotalAmount(&order); err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER999", StatusCode: 500}, map[string]interface{}{"function": "CreateOrders"})
	}

	if err := uc.repo.CreateOrders(c.UserContext(), &order); err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER999", StatusCode: 500}, map[string]interface{}{"function": "CreateOrders"})
	}

	return adap_http.ResponseCreateOrder(c, order)
}

func (uc *Orders) GetOrders(c *fiber.Ctx) error {
	filter := map[string]interface{}{}
	if c.Params("id") != "" {
		filter["id"] = c.Params("id")
	}

	page := 1
	if c.Query("page") != "" {
		page, _ = strconv.Atoi(c.Query("page"))
	}

	size := 10
	if c.Query("size") != "" {
		size, _ = strconv.Atoi(c.Query("size"))
	}

	totalOrder, err := uc.repo.FindCount(c.UserContext(), filter)
	if err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER999", StatusCode: 500}, map[string]interface{}{"function": "GetOrders"})
	}

	orders, err := uc.repo.FindWithPagination(c.UserContext(), size, ((page - 1) * size), "created_at asc", filter)
	if err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER999", StatusCode: 500}, map[string]interface{}{"function": "GetOrders"})
	}
	return adap_http.ResponseGetOrder(c, orders, totalOrder)
}

func (uc *Orders) UpdateOrder(c *fiber.Ctx) error {
	orderId := c.Params("id")
	var orderUpdateRequest entities.OrderUpdate
	if err := c.BodyParser(&orderUpdateRequest); err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "UpdateOrder"})
	}

	// validate request body
	err := uc.validator.Struct(orderUpdateRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER400", StatusCode: 400}, map[string]interface{}{"function": "UpdateOrder"})
		}
	}

	filter := map[string]interface{}{"id": orderId}
	if err := uc.repo.UpdateStatusOrder(c.UserContext(), filter, orderUpdateRequest.Status); err != nil {
		return adap_http.Response(c, entities.Response{Status: "ER", ErrorMessage: err.Error(), ErrorCode: "ER999", StatusCode: 500}, map[string]interface{}{"function": "UpdateOrder"})
	}

	return adap_http.ResponseUpdateOrder(c, orderId, orderUpdateRequest.Status)
}
