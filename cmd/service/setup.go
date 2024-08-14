package service

import (
	"chat-api/internal/api"
	"chat-api/internal/middleware"
	"chat-api/internal/repository"
	"chat-api/internal/service"
	"chat-api/internal/websocket"
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SetupMongoDB(mongoURI string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SetupHandlers(db *mongo.Database, jwtSecret string) (*api.UserHandler, *api.ChatHandler, *api.WebSocketHandler) {
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)

	userService := service.NewUserService(userRepo, []byte(jwtSecret))
	chatService := service.NewChatService(chatRepo)

	userHandler := api.NewUserHandler(userService)
	chatHandler := api.NewChatHandler(chatService)

	hub := websocket.NewHub()
	go hub.Run()
	wsHandler := api.NewWebSocketHandler(hub, chatRepo, []byte(jwtSecret))

	return userHandler, chatHandler, wsHandler
}

func SetupRouter(userHandler *api.UserHandler, chatHandler *api.ChatHandler, wsHandler *api.WebSocketHandler, jwtSecret string) *mux.Router {
	r := mux.NewRouter()
	jwtMiddleware := middleware.NewJWTMiddleware([]byte(jwtSecret))

	r.HandleFunc("/api/register", userHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/api/login", userHandler.LoginUser).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(jwtMiddleware.Middleware)

	protected.HandleFunc("/logout/{id}", userHandler.LogoutUser).Methods("POST")
	protected.HandleFunc("/send-message", chatHandler.SendMessage).Methods("POST")
	protected.HandleFunc("/chat-history/{senderID}/{receiverID}", chatHandler.GetChatHistory).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	protected.HandleFunc("/users/{id}/others", userHandler.GetOtherUsers).Methods("GET")

	r.HandleFunc("/ws", wsHandler.ServeWebSocket)

	return r
}

func StartServer(r *mux.Router, address string) {
	server := &http.Server{
		Handler:      r,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server started at", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", address, err)
		}
	}()

	serverShutdown(server)
}

func serverShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
