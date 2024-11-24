package controller

import (
	"merchant-bank-api/models"
	"merchant-bank-api/service"

	"github.com/gin-gonic/gin"

	"net/http"
)

type paymentController struct {
	service service.PaymentService
	rg      *gin.RouterGroup
}

func (c *paymentController) postPaymentHandlers(ctx *gin.Context) {
	var payload models.PaymentRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	data, err := c.service.PostPayment(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "filed to create payment")
		return
	}
	ctx.JSON(http.StatusOK, data)
}
func (c *paymentController) Route() {
	router := c.rg.Group("payment-merchant")
	router.POST("/", c.postPaymentHandlers)
}

func NewPaymentController(ps service.PaymentService, rg *gin.RouterGroup) *paymentController {
	return &paymentController{service: ps, rg: rg}
}
