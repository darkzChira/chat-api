package main

import (
	"chat-api/cmd/service"
	"chat-api/internal/config"
	"context"
	"log"
)

func main() {
	conf := config.LoadConfig()
	client, err := service.SetupMongoDB(conf.MongoURI)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(conf.DatabaseName)
	userHandler, chatHandler, wsHandler := service.SetupHandlers(db, conf.JWTSecret)
	r := service.SetupRouter(userHandler, chatHandler, wsHandler, conf.JWTSecret)

	service.StartServer(r, conf.ServerAddress)
}
