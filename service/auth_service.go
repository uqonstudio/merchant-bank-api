package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"merchant-bank-api/models"
	"merchant-bank-api/models/dto"
	"merchant-bank-api/util"
	"os"
	"time"
)

type AuthService interface {
	PostLogin(payload dto.LoginRequest) (dto.LoginResponse, error)
}

type authService struct {
	jwtservice JwtService
	cs         CustomerService
}

func (s *authService) PostLogin(payload dto.LoginRequest) (dto.LoginResponse, error) {
	// fmt.Println("payload : ", payload)
	customers, err := s.cs.GetAllCustomer()
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// fmt.Println("customers : ", customers)
	for i, customer := range customers {
		// fmt.Println("is user same : ", payload.Username, customer.Username)
		if s.isPasswordValid(payload.Password, customer.Password) && (payload.Username == customer.Username) {
			// fmt.Println("is selected user : ", customer)
			customers[i].LoggedIn = true
			// logHistory(customer.ID, "login")
			err := s.cs.UpdateCustomerLoggedInStatus(customer.Username, true)
			if err != nil {
				return dto.LoginResponse{}, err
			}
			return s.createLoginResponse(customer)
		}
	}

	return dto.LoginResponse{}, errors.New("invalid credentials")
}

func NewAuthService(jwtservice JwtService, cs CustomerService) AuthService {
	return &authService{jwtservice, cs}
}

func (s *authService) isPasswordValid(inputPassword, storedPassword string) bool {
	err := util.ComparePassword(inputPassword, storedPassword)
	return err == nil
}

func (s *authService) createLoginResponse(customer models.Customer) (dto.LoginResponse, error) {
	token, err := s.jwtservice.GenerateToken(customer)
	if err != nil {
		return dto.LoginResponse{}, errors.New("failed to generate token")
	}
	return token, nil
}

func logHistory(customerID, action string) {
	history := models.History{
		CustomerID: customerID,
		Action:     action,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Open the history file
	file, err := os.Open("database/history.json")
	if err != nil {
		log.Printf("Error opening history file: %v", err)
		return
	}
	defer file.Close()

	// Read existing history entries
	var histories []models.History
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&histories); err != nil && err != io.EOF {
		log.Printf("Error decoding history file: %v", err)
		return
	}

	// Append the new history entry
	histories = append(histories, history)

	// Write back to the file
	UpdateHistory(histories)
}

func UpdateHistory(histories []models.History) {
	// Open the history file with write permissions
	file, err := os.OpenFile("database/history.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Error opening history file for writing: %v", err)
		return
	}
	defer file.Close()

	// Write back to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(histories); err != nil {
		log.Printf("Error encoding history data: %v", err)
		return
	}
}
