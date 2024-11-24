package main

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func (s *Server) initialRoute() {
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func (s *Server) Start() {
	s.engine.Run(":8080")
}

func NewServer() *Server {
	return &Server{
		engine: gin.Default(),
	}
}
