package repository

import (
	"chat-app/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{
		collection: db.Collection("messages"),
	}
}

func (r *ChatRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *ChatRepository) GetMessagesBetweenUsers(ctx context.Context, senderID, receiverID string, limit int64) ([]models.Message, error) {
	messages := []models.Message{}

	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": senderID, "receiver_id": receiverID},
			{"sender_id": receiverID, "receiver_id": senderID},
		},
	}
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
