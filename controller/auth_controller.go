package controller

import (
	"merchant-bank-api/models/dto"
	"merchant-bank-api/service"

	"github.com/gin-gonic/gin"

	"net/http"
)

type authController struct {
	service service.AuthService
	rg      *gin.RouterGroup
}

func (c *authController) logoutHandler(ctx *gin.Context) {
	var logoutRequest dto.LogoutRequest
	if err := ctx.ShouldBindJSON(&logoutRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	_, err := c.service.Logout(logoutRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "filed to logout")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (c *authController) loginHandler(ctx *gin.Context) {
	var payload dto.LoginRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	data, err := c.service.PostLogin(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "filed to login")
		return
	}
	ctx.JSON(http.StatusOK, data)
}
func (c *authController) Route() {
	router := c.rg.Group("auth")
	router.POST("/login", c.loginHandler)
	router.POST("/logout", c.logoutHandler)
}

func NewAuthController(as service.AuthService, rg *gin.RouterGroup) *authController {
	return &authController{service: as, rg: rg}
}
