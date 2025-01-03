package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"merchant-bank-api/models"
	"os"
	"time"
)

// PaymentService defines the interface for payment operations.
// It includes a method to process payment requests.
type PaymentService interface {
	PostPayment(models.PaymentRequest) (models.Payment, error)
}

// paymentService is a concrete implementation of PaymentService.
// It handles payment processing and interacts with customer and history services.
type paymentService struct {
	cs CustomerService
	hs HistoryService
}

// PostPayment processes a payment request.
// It retrieves the logged-in customer, verifies the transaction, creates a payment record,
// and logs the payment history. Returns the created payment or an error.
func (s *paymentService) PostPayment(paymentRequest models.PaymentRequest) (models.Payment, error) {
	customer, err := s.getLoggedInCustomer(paymentRequest.CustomerID)
	if err != nil {
		fmt.Println("getLoggedInCustomer error: ", err)
		return models.Payment{}, err
	}

	if err := s.verifyTransaction(paymentRequest.TransactionID); err != nil {
		fmt.Println("verifyTransaction error: ", err)
		return models.Payment{}, err
	}

	payment, err := s.createPaymentRecord(customer, paymentRequest)
	if err != nil {
		fmt.Println("createPaymentRecord error: ", err)
		return models.Payment{}, err
	}

	if err := s.hs.LogHistory(customer.ID, "payment"); err != nil {
		fmt.Println("LogHistory error: ", err)
		return models.Payment{}, err
	}

	return payment, nil
}

// NewPaymentService creates a new instance of paymentService.
// It requires a CustomerService and a HistoryService to function.
func NewPaymentService(cs CustomerService, hs HistoryService) PaymentService {
	return &paymentService{cs, hs}
}

// getLoggedInCustomer retrieves a logged-in customer by ID.
// It returns the customer if found and logged in, otherwise returns an error.
func (s *paymentService) getLoggedInCustomer(customerID string) (*models.Customer, error) {
	customers, err := s.cs.GetAllCustomer()
	if err != nil {
		return nil, err
	}

	for _, customer := range customers {
		if customer.ID == customerID && customer.LoggedIn {
			return &customer, nil
		}
	}
	return nil, errors.New("unauthorized or invalid customer")
}

// verifyTransaction verifies the transaction ID.
// This function should implement logic to verify the transaction with a third-party service.
func (s *paymentService) verifyTransaction(transactionID string) error {
	// Implement logic to verify the transaction ID with the third-party service
	log.Printf("Verifying transaction ID: %s", transactionID)
	return nil
}

// createPaymentRecord creates and saves a new payment record.
// It returns the created payment or an error if the operation fails.
func (s *paymentService) createPaymentRecord(customer *models.Customer, paymentRequest models.PaymentRequest) (models.Payment, error) {
	payment := models.Payment{
		CustomerID:    customer.ID,
		MerchantID:    paymentRequest.MerchantID,
		Amount:        paymentRequest.Amount,
		TransactionID: paymentRequest.TransactionID,
		Timestamp:     time.Now().Format(time.RFC3339),
	}

	payments, err := s.loadPayments()
	if err != nil {
		return models.Payment{}, fmt.Errorf("failed to load payments: %v", err)
	}

	payments = append(payments, payment)

	if err := s.savePayments(payments); err != nil {
		return models.Payment{}, fmt.Errorf("failed to save payment: %v", err)
	}

	return payment, nil
}

// loadPayments loads payment data from a JSON file.
// It returns a slice of payments or an error if the operation fails.
func (s *paymentService) loadPayments() ([]models.Payment, error) {
	file, err := os.Open("database/payment.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Payment{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var payments []models.Payment
	if err := json.NewDecoder(file).Decode(&payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// savePayments saves payment data to a JSON file.
// It returns an error if the operation fails.
func (s *paymentService) savePayments(payments []models.Payment) error {
	file, err := os.OpenFile("database/payment.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(payments); err != nil {
		return err
	}

	return nil
}
