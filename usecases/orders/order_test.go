package usecases_orders

import (
	"bytes"
	"context"
	"encoding/json"
	"go-oms/entities"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateOrders(ctx context.Context, order *entities.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *mockRepository) FindWithPagination(ctx context.Context, limit, offset int, orderBy string, filter map[string]interface{}) ([]entities.Order, error) {
	args := m.Called(ctx, limit, offset, orderBy, filter)
	return args.Get(0).([]entities.Order), args.Error(1)
}

func (m *mockRepository) FindCount(ctx context.Context, filter map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepository) UpdateStatusOrder(ctx context.Context, filter map[string]interface{}, status string) error {
	args := m.Called(ctx, filter, status)
	return args.Error(0)
}

func TestCreateOrders(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mockRepository)
	validate := validator.New()
	uc := NewOrders(mockRepo, validate)

	// Setup mock
	mockRepo.On("CreateOrders", mock.Anything, mock.Anything).Return(nil)
	app.Post("/orders", uc.CreateOrders)
	tests := []struct {
		description  string
		requestBody  entities.Order
		expectStatus int
	}{
		{
			description:  "valid input",
			requestBody:  entities.Order{CustomerName: "Test 1", OrderItems: []entities.OrderItem{{ProductName: "product A", Quantity: 1, Price: 100}}},
			expectStatus: 200,
		},
		{
			description:  "Invalid input",
			requestBody:  entities.Order{OrderItems: []entities.OrderItem{{ProductName: "product A", Quantity: 1, Price: 100}}},
			expectStatus: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest("POST", "/orders/", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, test.expectStatus, resp.StatusCode)
		})
	}
}

func TestGetOrders(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mockRepository)
	validate := validator.New()
	uc := NewOrders(mockRepo, validate)

	// Setup mock
	mockRepo.On("FindCount", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockRepo.On("FindWithPagination", mock.Anything, 10, 0, "created_at asc", mock.Anything).
		Return([]entities.Order{
			{ID: 1, CustomerName: "Test"},
		}, nil)
	app.Get("/orders", uc.GetOrders)
	tests := []struct {
		description  string
		requestBody  map[string]interface{}
		expectStatus int
	}{
		{
			description:  "valid input",
			requestBody:  nil,
			expectStatus: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest("GET", "/orders?page=1&size=10", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, test.expectStatus, resp.StatusCode)

			// อ่าน body response
			defer resp.Body.Close()
			var result entities.ResponseGetOrder
			err := json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err)
			assert.IsType(t, entities.ResponseGetOrder{}, result)
		})
	}
}

func TestUpdateOrders(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mockRepository)
	validate := validator.New()
	uc := NewOrders(mockRepo, validate)

	// Setup mock
	mockRepo.On("UpdateStatusOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	app.Put("/orders/:id", uc.UpdateOrder)
	tests := []struct {
		description  string
		requestBody  map[string]interface{}
		expectStatus int
	}{
		{
			description: "valid input -> status completed",
			requestBody: map[string]interface{}{
				"status": "completed",
			},
			expectStatus: 200,
		},
		{
			description: "valid input -> status canceled",
			requestBody: map[string]interface{}{
				"status": "canceled",
			},
			expectStatus: 200,
		},
		{
			description: "invalid input -> status weird",
			requestBody: map[string]interface{}{
				"status": "weird",
			},
			expectStatus: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest("PUT", "/orders/1", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, test.expectStatus, resp.StatusCode)

			// อ่าน body response
			defer resp.Body.Close()
			var result entities.ResponseGetOrder
			err := json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err)
			assert.IsType(t, entities.ResponseGetOrder{}, result)
		})
	}
}
