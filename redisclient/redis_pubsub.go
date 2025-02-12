package redisclient

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/turgaysozen/algotrading/models"
	"github.com/turgaysozen/algotrading/services"
)

var ctx = context.Background()
var connected bool = false

func NewRedisClient() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	addr := redisHost + ":" + redisPort
	maxRetries := 5
	retryCount := 0

	var client *redis.Client

	for {
		client = redis.NewClient(&redis.Options{
			Addr: addr,
		})

		_, err = client.Ping(context.Background()).Result()
		if err != nil {
			connected = false
			retryCount++
			log.Printf("Error connecting to Redis, retrying %d/%d... %v", retryCount, maxRetries, err)
			client.Close()
			if retryCount >= maxRetries {
				log.Fatal("Max retry attempts reached, giving up.")
			}
			time.Sleep(5 * time.Second)
			continue
		}

		if !connected {
			log.Println("Connected to Redis")
			connected = true
		}

		break
	}

	return client
}

func Publish(channel string, message interface{}) {
	client := NewRedisClient()
	defer client.Close()

	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error serializing object:", err)
		return
	}

	err = client.Publish(ctx, channel, data).Err()
	if err != nil {
		log.Println("Error publishing to Redis:", err)
	}
}

func Subscribe() {
	client := NewRedisClient()
	sub := client.Subscribe(ctx, "order_book")

	ch := sub.Channel()
	for msg := range ch {
		var orderBook models.OrderBook
		err := json.Unmarshal([]byte(msg.Payload), &orderBook)
		if err != nil {
			log.Println("Error unmarshalling Redis message:", err)
			continue
		}

		go services.ProcessOrderBook(orderBook)
	}
}
