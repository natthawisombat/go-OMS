package entities

type Response struct {
	Status          string      `json:"status"`
	StatusCode      int         `json:"statusCode,omitempty"`
	Message         string      `json:"message,omitempty"`
	ErrorMessage    string      `json:"errorMessage,omitempty"`
	ErrorCode       string      `json:"errorCode,omitempty"`
	TransactionCode string      `json:"transactionCode,omitempty"`
	Data            interface{} `json:"data,omitempty"`
}

type ResponseOrder struct {
	OrderID      uint           `json:"order_id"`
	CustomerName string         `json:"customer_name"`
	TotalAmount  float64        `json:"total_amount"`
	Items        []ResponseItem ` json:"items"`
}

type ResponseItem struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type ResponseGetOrder struct {
	Total  int     `json:"total"`
	Orders []Order `json:"order"`
}

type ResponseUpdateOrder struct {
	OrderId int    `json:"order_id"`
	Status  string `json:"status"`
}
