package main

import (
	"merchant-bank-api/config"
	"merchant-bank-api/controller"
	"merchant-bank-api/middleware"
	"merchant-bank-api/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	am     middleware.AuthMiddleware
	ps     service.PaymentService
	as     service.AuthService
	cs     service.CustomerService
	js     service.JwtService
	engine *gin.Engine
}

func main() {
	NewServer().Start()
}

func (s *Server) initialRoute() {
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	routerGroup := s.engine.Group("/api")
	controller.NewCustomerController(s.cs, routerGroup).Route()      //get, post customer
	controller.NewAuthController(s.as, routerGroup).Route()          //auth/login, logout
	controller.NewPaymentController(s.ps, s.am, routerGroup).Route() //payment with middleware
}

func (s *Server) Start() {
	s.initialRoute()
	s.engine.Run(":8080")
}

func NewServer() *Server {
	c, _ := config.NewConfig()
	cService := service.NewCustomerService()
	jwtService := service.NewJwtService(c.JwtConfig)
	hService := service.NewHistoryService()
	aService := service.NewAuthService(jwtService, cService, hService)
	pService := service.NewPaymentService(cService, hService)
	authMidleware := middleware.NewAuthMiddleware(jwtService)

	return &Server{
		am:     authMidleware,
		ps:     pService,
		as:     aService,
		cs:     cService,
		js:     jwtService,
		engine: gin.Default(),
	}
}
