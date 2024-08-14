package websocket

import (
	"chat-api/internal/models"
	"chat-api/internal/repository"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Connection struct {
	Ws      *websocket.Conn
	Send    chan []byte
	MsgRepo *repository.ChatRepository
}

func (c *Connection) ReadPump(hub *Hub) {

	var userID string

	defer func() {
		// Notify all users that this user has gone offline
		if userID != "" {
			userOfflineNotification := map[string]string{
				"type":    "user_offline",
				"user_id": userID,
			}
			notificationMessage, _ := json.Marshal(userOfflineNotification)
			hub.BroadcastMessage(notificationMessage)
		}

		hub.Unregister <- c
		c.Ws.Close()
	}()

	for {
		_, messageData, err := c.Ws.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var msg models.Message
		err = json.Unmarshal(messageData, &msg)
		if err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		msg.Timestamp = time.Now()

		err = c.MsgRepo.SaveMessage(context.Background(), &msg)
		if err != nil {
			log.Println("failed to save message:", err)
			continue
		}

		hub.Broadcast <- messageData
	}
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
