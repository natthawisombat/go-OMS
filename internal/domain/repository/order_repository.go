package repository

import (
	"context"

	"go-oms/internal/domain/entities"
)

type Repository interface {
	CreateOrders(ctx context.Context, order *entities.Order) error
	FindWithPagination(ctx context.Context, limit int, offset int, orderBy string, filter map[string]interface{}) ([]entities.Order, error)
	FindCount(ctx context.Context, filter map[string]interface{}) (int64, error)
	UpdateStatusOrder(ctx context.Context, filter map[string]interface{}, status string) error
}
