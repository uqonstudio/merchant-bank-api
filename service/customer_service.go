package service

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"merchant-bank-api/models"
	"merchant-bank-api/models/dto"
)

type CustomerService interface {
	GetAllCustomer() ([]models.Customer, error)
	PostCustomer(payload dto.CustomerPayload) (models.Customer, error)
}

type customerService struct{}

func (s *customerService) GetAllCustomer() ([]models.Customer, error) {
	// Open the customers file
	file, err := os.Open("database/customer.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the customers
	var customers []models.Customer
	if err := json.NewDecoder(file).Decode(&customers); err != nil {
		return nil, err
	}

	// Return the list of customers
	return customers, nil
}

func (s *customerService) PostCustomer(payload dto.CustomerPayload) (models.Customer, error) {
	// Read existing customers
	file, err := os.Open("database/customer.json")
	if err != nil {
		return models.Customer{}, err
	}
	defer file.Close()

	var customers []models.Customer
	json.NewDecoder(file).Decode(&customers)

	// Hash the password
	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		return models.Customer{}, err
	}

	// Create a new customer record
	newCustomer := models.Customer{
		ID:       fmt.Sprintf("%d", len(customers)+1),
		Username: payload.Username,
		Password: hashedPassword,
		LoggedIn: false,
	}

	// Append the new customer
	customers = append(customers, newCustomer)

	// Write back to the file
	fileData, err := json.MarshalIndent(customers, "", "  ")
	if err != nil {
		return models.Customer{}, err
	}

	err = os.WriteFile("database/customer.json", fileData, 0644)
	if err != nil {
		return models.Customer{}, err
	}

	return newCustomer, nil
}

func NewCustomerService() CustomerService {
	return &customerService{}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
