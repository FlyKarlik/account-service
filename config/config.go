package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config of service
type Config struct {
	ServerHost string
	NatsHost   string
	DbDsn      string
	JwtSecret  string
	AesSecret  string
	SentryDSN  string
	JaegerHost string
}

// NewConfig generate new config
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("NewConfig godotenv load failed")
	}

	return &Config{
		ServerHost: os.Getenv("SERVER_HOST"),
		NatsHost:   os.Getenv("NATS_HOST"),
		DbDsn:      os.Getenv("DB_DSN"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
		AesSecret:  os.Getenv("AES_SECRET"),
		SentryDSN:  os.Getenv("SENTRY_DSN"),
		JaegerHost: os.Getenv("JAEGER_HOST"),
	}, nil
}
