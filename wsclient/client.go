package wsclient

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/turgaysozen/algotrading/models"
	"github.com/turgaysozen/algotrading/monitoring/metrics"
	"github.com/turgaysozen/algotrading/redisclient"
)

var Connected bool = false

func ConnectWebSocket() (*websocket.Conn, error) {
	url := os.Getenv("WEB_SOCKET_URL")

	if url == "" {
		log.Fatal("WebSocket URL not set in .env file")
		metrics.RecordError("websocket_url_missing")
	}

	var conn *websocket.Conn
	var connectErr error
	maxRetries := 5
	retryCount := 0

	for {
		conn, _, connectErr = websocket.DefaultDialer.Dial(url, nil)
		if connectErr != nil {
			Connected = false
			retryCount++
			log.Printf("Error connecting to WebSocket, retrying %d/%d... %v", retryCount, maxRetries, connectErr)
			metrics.RecordError("websocket_connection_error")
			if retryCount >= maxRetries {
				log.Fatal("Max retry attempts reached, giving up.")
				metrics.RecordDataLoss("websocket_connection_max_retries")
				return nil, connectErr
			}
			time.Sleep(5 * time.Second)
			continue
		}

		Connected = true
		log.Println("WebSocket connected to:", url)
		break
	}

	Connected = true

	return conn, nil
}

func ProcessWebSocketMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			Connected = false
			log.Println("Error reading WebSocket message:", err)
			metrics.RecordError("websocket_connection_error")
			conn.Close()
			conn, err = ConnectWebSocket()
			if err != nil {
				log.Fatal("Error reconnecting to WebSocket:", err)
			}
			continue
		}

		Connected = true

		// track latency for orderbook avg processing
		metrics.SetStartTime("orderbook_avg")

		var orderBook models.OrderBook
		err = json.Unmarshal(msg, &orderBook)
		if err != nil {
			log.Println("Error unmarshalling WebSocket message:", err)
			metrics.RecordError("json_unmarshal_error")
			metrics.RecordDataLoss("json_unmarshal_data_loss")
			continue
		}

		redisclient.Publish("order_book", orderBook)
	}
}
