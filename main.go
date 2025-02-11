package main

import (
	"log"

	"github.com/turgaysozen/algotrading/redisclient"
	"github.com/turgaysozen/algotrading/wsclient"
)

func main() {
	conn, err := wsclient.ConnectWebSocket()
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	go wsclient.ProcessWebSocketMessages(conn)

	go redisclient.Subscribe()

	select {}
}
