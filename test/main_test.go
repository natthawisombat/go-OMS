package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-oms/configs"
	"go-oms/entities"
	pkglogger "go-oms/pkg/logger"
	"go-oms/routers"
	"net/http/httptest"
	"os/signal"
	"syscall"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	godotenv.Load("../.env")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	logger := pkglogger.NewLoggerFromEnv()
	cfg := configs.NewApp(logger)
	if err := cfg.SetApp(ctx); err != nil {
		fmt.Println(err)
		logger.Fatal(err)
	}
	routers.SetupRoutes(cfg)
	return cfg.App
}

func TestCreateOrder(t *testing.T) {
	app := setupTestApp()
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
			resp, _ := app.Test(req)

			assert.Equal(t, test.expectStatus, resp.StatusCode)
		})
	}
}

func TestGetOrders(t *testing.T) {
	app := setupTestApp()
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

func TestUpdateOrder(t *testing.T) {
	app := setupTestApp()
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
