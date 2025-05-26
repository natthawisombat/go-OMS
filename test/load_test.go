package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

type OrderItem struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type OrderRequest struct {
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
}

func TestCreateOrdersConcurrently(t *testing.T) {
	const total = 1000
	var wg sync.WaitGroup
	wg.Add(total)

	transport := &http.Transport{
		MaxIdleConns:    200,
		MaxConnsPerHost: 200,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   100 * time.Second,
	}

	errCh := make(chan error, total) // ช่องเก็บ error
	start := time.Now()

	for i := 0; i < total; i++ {
		go func(i int) {
			defer wg.Done()

			order := OrderRequest{
				CustomerName: "User " + strconv.Itoa(i),
				Items: []OrderItem{
					{ProductName: "Product A", Quantity: 1, Price: 99.99},
				},
			}

			body, _ := json.Marshal(order)
			resp, err := client.Post("http://localhost:8080/orders/", "application/json", bytes.NewReader(body))
			if err != nil {
				errCh <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				errCh <- &StatusCodeError{StatusCode: resp.StatusCode}
				return
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	elapsed := time.Since(start)
	t.Logf("✅ Completed %d requests in %s", total, elapsed)

	// ตรวจว่ามี error ไหม
	if len(errCh) > 0 {
		for err := range errCh {
			t.Errorf("❌ Request failed: %v", err)
		}
		t.Fatalf("⛔ Total failed: %d/%d", len(errCh), total)
	}
}

type StatusCodeError struct {
	StatusCode int
}

func (e *StatusCodeError) Error() string {
	return "invalid status code: " + strconv.Itoa(e.StatusCode)
}
