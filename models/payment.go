package models

type PaymentRequest struct {
	TransactionID string  `json:"transaction_id"`
	CustomerID    string  `json:"customer_id"`
	MerchantID    string  `json:"merchant_id"`
	Amount        float64 `json:"amount"`
}

type Payment struct {
	TransactionID string  `json:"transaction_id"`
	CustomerID    string  `json:"customer_id"`
	MerchantID    string  `json:"merchant_id"`
	Amount        float64 `json:"amount"`
	Timestamp     string  `json:"timestamp"`
}
