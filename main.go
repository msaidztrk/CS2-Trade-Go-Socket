package main

import (
	"log"

	"CS2-Trade-Go-Socket/auth"
	"CS2-Trade-Go-Socket/bidding"
	"CS2-Trade-Go-Socket/websocket"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Authenticate and get WebSocket token
	token, err := auth.Authenticate()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Connect to WebSocket
	conn, err := websocket.Connect("wss://ws.csgoempire.com", map[string]string{
		"Authorization": "Bearer " + token,
	})
	if err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}
	defer conn.Close()

	// Subscribe to the bidding channel
	if err := websocket.Subscribe(conn, "bidding"); err != nil {
		log.Fatalf("Failed to subscribe to bidding channel: %v", err)
	}

	// Handle incoming messages
	websocket.ReadMessages(conn, func(data map[string]interface{}) {
		if data["channel"] == "bidding" {
			bidding.HandleBiddingUpdate(conn, data)
		}
	})
}
