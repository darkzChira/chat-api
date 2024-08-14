package service

import (
	"chat-app/internal/models"
	"chat-app/internal/repository"
	"context"
	"time"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{chatRepo: chatRepo}
}

func (s *ChatService) SendMessage(ctx context.Context, senderID, receiverID, content string) error {
	message := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Timestamp:  time.Now(),
	}
	return s.chatRepo.SaveMessage(ctx, message)
}

func (s *ChatService) GetChatHistory(ctx context.Context, senderID, receiverID string, limit int64) ([]models.Message, error) {
	return s.chatRepo.GetMessagesBetweenUsers(ctx, senderID, receiverID, limit)
}
