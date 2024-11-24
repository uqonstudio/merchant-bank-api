// models/customer.go
package models

type Customer struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	LoggedIn bool   `json:"logged_in"`
}
