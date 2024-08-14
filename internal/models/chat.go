package models

import "time"

type Message struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	SenderID   string    `json:"sender_id" bson:"sender_id"`
	ReceiverID string    `json:"receiver_id" bson:"receiver_id"`
	Content    string    `json:"content" bson:"content"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
}
