package usecases_orders

import "go-oms/entities"

func calcurateTotalAmount(order *entities.Order) error {
	for _, v := range order.OrderItems {
		order.TotalAmount += float64(v.Quantity) * v.Price
	}

	return nil
}
