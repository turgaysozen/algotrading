package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/turgaysozen/algotrading/db"
	"github.com/turgaysozen/algotrading/monitoring"
	"github.com/turgaysozen/algotrading/monitoring/metrics"
	"github.com/turgaysozen/algotrading/redisclient"
	"github.com/turgaysozen/algotrading/wsclient"
)

func main() {
	_, err := db.InitializeDB()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	redisclient.InitRedisClient()

	conn, err := wsclient.ConnectWebSocket()
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/healthz", monitoring.LivenessHandler)
		http.HandleFunc("/readiness", monitoring.ReadinessHandler)

		log.Println("Client metrics and health checks available at http://localhost:8080/metrics, /healthz, /readiness")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go func() {
		for {
			metrics.CollectSystemMetrics()
			time.Sleep(10 * time.Second)
		}
	}()

	go wsclient.ProcessWebSocketMessages(conn)

	go redisclient.Subscribe()

	select {}
}
