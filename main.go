package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/turgaysozen/algotrading/db"
	"github.com/turgaysozen/algotrading/metrics"
	"github.com/turgaysozen/algotrading/redisclient"
	"github.com/turgaysozen/algotrading/wsclient"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Client metrics available at http://localhost:8080/metrics")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go func() {
		for {
			metrics.CollectSystemMetrics()
			time.Sleep(10 * time.Second)
		}
	}()

	db, err := db.InitializeDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	conn, err := wsclient.ConnectWebSocket()
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	go wsclient.ProcessWebSocketMessages(conn)

	go redisclient.Subscribe()

	select {}
}
