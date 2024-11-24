package controller

import (
	"fmt"
	"merchant-bank-api/models/dto"
	"merchant-bank-api/service"

	"github.com/gin-gonic/gin"

	"net/http"
)

type customerController struct {
	service service.CustomerService
	rg      *gin.RouterGroup
}

func (c *customerController) getAllHandlers(ctx *gin.Context) {
	data, _ := c.service.GetAllCustomer()
	ctx.JSON(http.StatusOK, data)
}

// postHandler handles POST requests to create a new customer.
// It expects a JSON payload containing the customer details in the request body.
// If the payload is invalid, it returns a 200 status code with an error message in the response body.
// If the customer is successfully created, it returns a 200 status code with the created customer data in the response body.
// If an error occurs during the creation process, it returns a 500 status code with a generic error message.
func (c *customerController) postHandler(ctx *gin.Context) {
	var payload dto.CustomerPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	fmt.Println("payload : ", payload)
	data, err := c.service.PostCustomer(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "filed to create user")
		return
	}
	ctx.JSON(http.StatusOK, data)
}
func (c *customerController) Route() {
	router := c.rg.Group("customers")
	router.GET("/", c.getAllHandlers)
	router.POST("/", c.postHandler)
}

func NewCustomerController(cs service.CustomerService, rg *gin.RouterGroup) *customerController {
	return &customerController{service: cs, rg: rg}
}
