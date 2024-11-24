package service

import (
	"errors"
	"merchant-bank-api/models"
	"merchant-bank-api/models/dto"
	"merchant-bank-api/util"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	// PostLogin handles user login requests.
	// It takes a LoginRequest payload and returns a LoginResponse or an error.
	PostLogin(payload dto.LoginRequest) (dto.LoginResponse, error)
	// Logout handles user logout requests.
	// It takes a LogoutRequest payload and returns a message or an error.
	Logout(payload dto.LogoutRequest) (string, error)
}

// authService is a concrete implementation of the AuthService interface.
type authService struct {
	jwtservice JwtService
	cs         CustomerService
	hs         HistoryService
}

// Logout processes a logout request for a customer.
// It checks if the customer is logged in and updates their status.
// Returns a success message or an error if the operation fails.
func (s *authService) Logout(payload dto.LogoutRequest) (string, error) {
	customers, err := s.cs.GetAllCustomer()
	if err != nil {
		return "", err
	}

	for i, customer := range customers {
		if customer.ID == payload.CustomerID && customer.LoggedIn {
			return s.processLogout(&customers[i])
		}
	}
	return "Unauthorized or invalid customer", nil
}

// PostLogin processes a login request for a customer.
// It validates the credentials and updates the customer's logged-in status.
// Returns a LoginResponse with a token or an error if the credentials are invalid.
func (s *authService) PostLogin(payload dto.LoginRequest) (dto.LoginResponse, error) {
	customers, err := s.cs.GetAllCustomer()
	if err != nil {
		return dto.LoginResponse{}, err
	}

	for i, customer := range customers {
		if s.isPasswordValid(payload.Password, customer.Password) && (payload.Username == customer.Username) {
			customers[i].LoggedIn = true
			_ = s.hs.LogHistory(customer.ID, "loggin")
			err := s.cs.UpdateCustomerLoggedInStatus(customer.Username, true)
			if err != nil {
				return dto.LoginResponse{}, err
			}
			return s.createLoginResponse(customer)
		}
	}

	return dto.LoginResponse{}, errors.New("invalid credentials")
}

// NewAuthService creates a new instance of authService with the provided dependencies.
func NewAuthService(jwtservice JwtService, cs CustomerService, hs HistoryService) AuthService {
	return &authService{jwtservice, cs, hs}
}

// isPasswordValid checks if the provided password matches the stored password.
// Returns true if the passwords match, false otherwise.
func (s *authService) isPasswordValid(inputPassword, storedPassword string) bool {
	err := util.ComparePassword(inputPassword, storedPassword)
	return err == nil
}

// createLoginResponse generates a LoginResponse containing a JWT token for the customer.
// Returns the LoginResponse or an error if token generation fails.
func (s *authService) createLoginResponse(customer models.Customer) (dto.LoginResponse, error) {
	token, err := s.jwtservice.GenerateToken(customer)
	if err != nil {
		return dto.LoginResponse{}, errors.New("failed to generate token")
	}
	return token, nil
}

// processLogout updates the customer's logged-in status and logs the logout action.
// Returns a success message or an error if the operation fails.
func (s *authService) processLogout(customer *models.Customer) (string, error) {
	customer.LoggedIn = false
	if err := s.hs.LogHistory(customer.ID, "logout"); err != nil {
		return "", err
	}
	if err := s.cs.UpdateCustomerLoggedInStatus(customer.Username, false); err != nil {
		return "", err
	}
	return "Logout successful", nil
}
