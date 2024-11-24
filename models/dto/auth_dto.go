package dto

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type JwtCustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type CustomerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
