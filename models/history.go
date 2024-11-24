// models/history.go
package models

type History struct {
	CustomerID string `json:"customer_id"`
	Action     string `json:"action"`
	Timestamp  string `json:"timestamp"`
}
