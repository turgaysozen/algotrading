package wsclient

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/turgaysozen/algotrading/metrics"
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

	var conn *websocket.Conn
	var connectErr error
	maxRetries := 5
	retryCount := 0

	for {
		conn, _, connectErr = websocket.DefaultDialer.Dial(url, nil)
		if connectErr != nil {
			retryCount++
			log.Printf("Error connecting to WebSocket, retrying %d/%d... %v", retryCount, maxRetries, connectErr)
			if retryCount >= maxRetries {
				log.Fatal("Max retry attempts reached, giving up.")
				return nil, connectErr
			}
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("WebSocket connected to:", url)
		break
	}

	return conn, nil
}

func ProcessWebSocketMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			conn.Close()
			conn, err = ConnectWebSocket()
			if err != nil {
				log.Fatal("Error reconnecting to WebSocket:", err)
			}
			continue
		}

		// track latency for orderbook single data processing and avg processing
		latencyTrackingID := uuid.New().String()
		metrics.SetStartTime("orderbook_single", latencyTrackingID)
		metrics.SetStartTime("orderbook_avg", "")

		var orderBook models.OrderBook
		err = json.Unmarshal(msg, &orderBook)
		if err != nil {
			log.Println("Error unmarshalling WebSocket message:", err)
			continue
		}

		orderBook.LatencyTrackingID = latencyTrackingID

		redisclient.Publish("order_book", orderBook)
	}
}
