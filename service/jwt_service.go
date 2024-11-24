package service

import (
	"errors"
	"fmt"
	"merchant-bank-api/config"
	"merchant-bank-api/models"
	"merchant-bank-api/models/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JwtService defines the interface for JWT operations, including generating and verifying tokens.
type JwtService interface {
	// GenerateToken generates a JWT token for a given customer payload.
	// Returns a LoginResponse containing the token or an error if token generation fails.
	GenerateToken(payload models.Customer) (dto.LoginResponse, error)
	// VerificationToken verifies a given JWT token string.
	// Returns the token claims if valid, or an error if verification fails.
	VerificationToken(token string) (jwt.MapClaims, error)
}

// jwtService is a private struct that implements the JwtService interface.
type jwtService struct {
	conf config.JwtConfig // Configuration for JWT, including issuer and signing key.
}

// GenerateToken creates a JWT token using the customer payload.
// It sets custom claims including UserId and standard claims like Issuer, ExpiresAt, and IssuedAt.
func (js *jwtService) GenerateToken(payload models.Customer) (dto.LoginResponse, error) {
	fmt.Println("Generate token :", payload)
	claims := dto.JwtCustomClaims{
		UserId: payload.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    js.conf.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(js.conf.Key))
	if err != nil {
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{Token: ss}, nil
}

// VerificationToken parses and verifies a JWT token string.
// It checks the token's validity, issuer, and claims.
func (js *jwtService) VerificationToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.conf.Key), nil
	})
	if err != nil {
		return nil, errors.New("failed parse token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || claims["iss"] != js.conf.Issuer || !ok {
		return nil, errors.New("invalid issuer or claims")
	}
	return claims, nil
}

// NewJwtService creates a new instance of JwtService with the provided configuration.
func NewJwtService(conf config.JwtConfig) JwtService {
	return &jwtService{conf: conf}
}
