package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MongoURI      string
	DatabaseName  string
	ServerAddress string
	WSAddress     string
	JWTSecret     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		MongoURI:      os.Getenv("MONGO_URI"),
		DatabaseName:  os.Getenv("DATABASE_NAME"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		WSAddress:     os.Getenv("WS_ADDRESS"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}
}
