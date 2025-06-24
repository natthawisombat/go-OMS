package entities

import "time"

type Order struct {
	ID           uint        `gorm:"primaryKey,autoIncrement" json:"order_id"`
	CustomerName string      `gorm:"type:varchar(100)" validate:"required" json:"customer_name"`
	TotalAmount  float64     `gorm:"type:decimal(10,2)" json:"total_amount"`
	Status       string      `gorm:"type:varchar(20)" json:"status"`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey,autoIncrement" json:"item_id"`
	OrderID     uint    `gorm:"index" json:"order_id"`
	ProductName string  `gorm:"type:varchar(100)" json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `gorm:"type:decimal(10,2)" json:"price"`
}

type OrderUpdate struct {
	Status string `validate:"required,oneof='completed' 'canceled'" json:"status"`
}
