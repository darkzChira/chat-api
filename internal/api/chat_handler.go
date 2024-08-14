package api

import (
	"chat-app/internal/models"
	"chat-app/internal/service"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := h.chatService.SendMessage(context.Background(), message.SenderID, message.ReceiverID, message.Content)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderID := vars["senderID"]
	receiverID := vars["receiverID"]
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	messages, err := h.chatService.GetChatHistory(context.Background(), senderID, receiverID, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve chat history", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}
