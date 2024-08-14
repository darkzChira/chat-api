package api

import (
	"chat-app/internal/repository"
	ws "chat-app/internal/websocket"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketHandler struct {
	hub     *ws.Hub
	msgRepo *repository.ChatRepository
	secret  []byte
}

func NewWebSocketHandler(hub *ws.Hub, msgRepo *repository.ChatRepository, secret []byte) *WebSocketHandler {
	return &WebSocketHandler{hub: hub, msgRepo: msgRepo, secret: secret}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return h.secret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID := claims["user_id"].(string)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}

	client := &ws.Connection{
		Ws:      conn,
		Send:    make(chan []byte, 256),
		MsgRepo: h.msgRepo,
	}

	h.hub.Register <- client

	userOnlineNotification := map[string]string{
		"type":    "user_online",
		"user_id": userID,
	}
	notificationMessage, _ := json.Marshal(userOnlineNotification)
	h.hub.BroadcastMessage(notificationMessage)

	go client.WritePump()
	client.ReadPump(h.hub)
}
