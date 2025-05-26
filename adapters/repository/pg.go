package adap_repository

import (
	"context"
	"fmt"
	"go-oms/entities"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) CreateOrders(ctx context.Context, order *entities.Order) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx := repo.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	tx = tx.WithContext(ctx)

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil

}

func (repo *Repository) FindWithPagination(
	ctx context.Context,
	limit int,
	offset int,
	orderBy string,
	filter map[string]interface{},
) ([]entities.Order, error) {
	var orders []entities.Order
	tx := repo.db.WithContext(ctx).Model(&entities.Order{})

	for k, v := range filter {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := tx.Preload("OrderItems").Order(orderBy).Limit(limit).Offset(offset).Find(&orders).Error
	return orders, err
}

func (repo *Repository) FindCount(
	ctx context.Context,
	filter map[string]interface{},
) (int64, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(&entities.Order{})

	for k, v := range filter {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := tx.Count(&count).Error
	return count, err
}

func (repo *Repository) UpdateStatusOrder(ctx context.Context, filter map[string]interface{}, status string) error {
	tx := repo.db.WithContext(ctx).Model(&entities.Order{})
	for k, v := range filter {
		tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := tx.Update("status", status).Error
	return err
}
