package adap_repository

import (
	"context"
	"fmt"
	"testing"

	"go-oms/internal/domain/entities"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&entities.Order{}, &entities.OrderItem{})
	assert.NoError(t, err)

	return db
}

func TestCreateOrders(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	order := &entities.Order{
		CustomerName: "Test User",
		OrderItems: []entities.OrderItem{
			{ProductName: "item A", Quantity: 1, Price: 100},
		},
	}

	err := repo.CreateOrders(context.Background(), order)
	assert.NoError(t, err)

	var count int64
	err = db.Model(&entities.Order{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	db.Migrator().DropTable(&entities.Order{}, &entities.OrderItem{})
}

func TestFindWithPagination(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	// Prepare test data
	db.Create(&entities.Order{
		CustomerName: "Test",
		OrderItems: []entities.OrderItem{
			{ProductName: "item", Quantity: 1, Price: 10},
		},
	})

	// filter := map[string]interface{}{
	// 	"customer_name": "Test",
	// }

	orders, err := repo.FindWithPagination(context.Background(), 10, 0, "id desc", nil)
	fmt.Println(orders)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	db.Migrator().DropTable(&entities.Order{}, &entities.OrderItem{})
}

func TestFindCount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	db.Create(&entities.Order{CustomerName: "Test"})

	count, err := repo.FindCount(context.Background(), map[string]interface{}{
		"customer_name": "Test",
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	db.Migrator().DropTable(&entities.Order{}, &entities.OrderItem{})
}

func TestUpdateStatusOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	order := entities.Order{
		CustomerName: "Test",
		Status:       "pending",
	}
	db.Create(&order)

	filter := map[string]interface{}{"customer_name": "Test"}
	err := repo.UpdateStatusOrder(context.Background(), filter, "completed")
	assert.NoError(t, err)

	var updated entities.Order
	db.First(&updated, order.ID)
	assert.Equal(t, "completed", updated.Status)
	db.Migrator().DropTable(&entities.Order{}, &entities.OrderItem{})
}
