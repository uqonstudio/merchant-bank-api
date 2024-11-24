package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type JwtConfig struct {
	Key    string
	Durasi time.Duration
	Issuer string
}

type Config struct {
	JwtConfig
}

func (c *Config) readConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	longTime, _ := strconv.Atoi(os.Getenv("JWT_LIFE_TIME"))
	c.JwtConfig = JwtConfig{
		Key:    os.Getenv("JWT_KEY"),
		Durasi: time.Duration(longTime),
		Issuer: os.Getenv("JWT_ISSUER_NAME"),
	}

	if c.JwtConfig.Key == "" {
		return errors.New("JWT_KEY not set in.env")
	}
	return nil
}

func NewConfig() (*Config, error) {
	config := &Config{}
	err := config.readConfig()
	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}
