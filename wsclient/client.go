package wsclient

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/turgaysozen/algotrading/models"
	"github.com/turgaysozen/algotrading/redisclient"
)

func ConnectWebSocket() (*websocket.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("WEB_SOCKET_URL")

	if url == "" {
		log.Fatal("WebSocket URL not set in .env file")
	}

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	log.Println("WebSocket connected to:", url)
	return c, nil
}

func ProcessWebSocketMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading WebSocket message:", err)
		}

		var orderBook models.OrderBook
		err = json.Unmarshal(msg, &orderBook)
		if err != nil {
			log.Println("Error unmarshalling WebSocket message:", err)
			continue
		}

		redisclient.Publish("order_book", orderBook)
	}
}
