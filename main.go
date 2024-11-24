package main

import (
	"merchant-bank-api/config"
	"merchant-bank-api/controller"
	"merchant-bank-api/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
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
	controller.NewCustomerController(s.cs, routerGroup).Route()
	controller.NewAuthController(s.as, routerGroup).Route()
}

func (s *Server) Start() {
	s.initialRoute()
	s.engine.Run(":8080")
}

func NewServer() *Server {
	c, _ := config.NewConfig()
	cService := service.NewCustomerService()
	jwtService := service.NewJwtService(c.JwtConfig)
	aService := service.NewAuthService(jwtService, cService)

	return &Server{
		as:     aService,
		cs:     cService,
		js:     jwtService,
		engine: gin.Default(),
	}
}
