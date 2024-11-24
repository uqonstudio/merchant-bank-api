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

type JwtService interface {
	GenerateToken(payload models.Customer) (dto.LoginResponse, error)
	VerificationToken(token string) (jwt.MapClaims, error)
}

type jwtService struct {
	conf config.JwtConfig
}

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

func NewJwtService(conf config.JwtConfig) JwtService {
	return &jwtService{conf: conf}
}
