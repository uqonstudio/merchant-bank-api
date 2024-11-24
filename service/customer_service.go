package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"merchant-bank-api/models"
	"merchant-bank-api/models/dto"
)

type CustomerService interface {
	GetAllCustomer() ([]models.Customer, error)
	PostCustomer(payload dto.CustomerPayload) (models.Customer, error)
	UpdateCustomerLoggedInStatus(username string, status bool) error
	// LoadCustomers() ([]models.Customer, error)
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

func (s *customerService) UpdateCustomerLoggedInStatus(username string, status bool) error {
	customers, err := s.GetAllCustomer()
	if err != nil {
		return err
	}

	updated := false
	for i, customer := range customers {
		if customer.Username == username {
			// fmt.Println("customers update :", customer.Username, username)
			customers[i].LoggedIn = status
			updated = true
			break
		}
		// fmt.Println("customers nu :", customer.Username, username)
	}

	if !updated {
		return errors.New("customer not found")
	}

	if err := s.saveCustomers(customers); err != nil {
		return err
	}

	return nil
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

func (s *customerService) saveCustomers(customers []models.Customer) error {
	file, err := os.Create("database/customer.json")
	if err != nil {
		return errors.New("failed to open customer database for writing")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(customers); err != nil {
		log.Printf("Error encoding customers data: %v", err)
		return errors.New("failed to encode customer data")
	}

	return nil
}
