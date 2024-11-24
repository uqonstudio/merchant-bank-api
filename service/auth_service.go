package service

import (
	"errors"
	"merchant-bank-api/models"
	"merchant-bank-api/models/dto"
	"merchant-bank-api/util"
)

type AuthService interface {
	PostLogin(payload dto.LoginRequest) (dto.LoginResponse, error)
	Logout(payload dto.LogoutRequest) (string, error)
}

type authService struct {
	jwtservice JwtService
	cs         CustomerService
	hs         HistoryService
}

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

func NewAuthService(jwtservice JwtService, cs CustomerService, hs HistoryService) AuthService {
	return &authService{jwtservice, cs, hs}
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

// processLogout updates the customer's logged-in status and logs the logout action
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
