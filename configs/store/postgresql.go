package store

import (
	"fmt"
	"go-oms/entities"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresqlStore struct {
	Store *gorm.DB
}

func ConnectPostgresql() (*PostgresqlStore, error) {
	PG_PORT, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		os.Getenv("PG_HOST"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DBNAME"), PG_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)                 // max concurrent open connections
	sqlDB.SetMaxIdleConns(20)                 // max idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // connection lifetime
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // optional: idle time

	err = db.AutoMigrate(&entities.Order{}, &entities.OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &PostgresqlStore{Store: db}, nil
}
