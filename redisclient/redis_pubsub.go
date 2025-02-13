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
	"github.com/turgaysozen/algotrading/monitoring/metrics"
	"github.com/turgaysozen/algotrading/services"
)

var ctx = context.Background()
var connected bool = false
var redisClient *redis.Client

func NewRedisClient() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		metrics.RecordError("env_file_load_error")
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
			metrics.RecordError("redis_connection_error")
			client.Close()
			if retryCount >= maxRetries {
				log.Fatal("Max retry attempts reached, giving up.")
				metrics.RecordDataLoss("redis_connection_max_retries")
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

func InitRedisClient() {
	if redisClient == nil {
		redisClient = NewRedisClient()
	}
}

func RedisHealth() error {
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error pinging Redis:", err)
		return err
	}
	return nil
}

func Publish(channel string, message interface{}) {
	InitRedisClient()

	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Error serializing object:", err)
		metrics.RecordError("redis_serialization_error")
		return
	}

	err = redisClient.Publish(ctx, channel, data).Err()
	if err != nil {
		log.Println("Error publishing to Redis:", err)
		metrics.RecordError("redis_publish_error")
	}
}

func Subscribe() {
	InitRedisClient()

	sub := redisClient.Subscribe(ctx, "order_book")

	ch := sub.Channel()
	for msg := range ch {
		var orderBook models.OrderBook
		err := json.Unmarshal([]byte(msg.Payload), &orderBook)
		if err != nil {
			log.Println("Error unmarshalling Redis message:", err)
			metrics.RecordError("redis_unmarshal_error")
			continue
		}

		go services.ProcessOrderBook(orderBook)
	}
}
