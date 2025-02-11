package main

import (
	"log"

	"github.com/turgaysozen/algotrading/db"
	"github.com/turgaysozen/algotrading/redisclient"
	"github.com/turgaysozen/algotrading/wsclient"
)

func main() {
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
