package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const csgoEmpireWSURL = "wss://trade.csgoempire.com/s/?EIO=3&transport=websocket"

func main() {
	// Get API key from environment
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found: %v", err)
	}

	apiKey := os.Getenv("CSGO_EMPIRE_API_KEY")
	if apiKey == "" {
		log.Fatal("CSGO_EMPIRE_API_KEY not found in environment variables")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Create WebSocket connection with auth header
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+apiKey)

	conn, _, err := websocket.DefaultDialer.Dial(csgoEmpireWSURL, header)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to CSGOEmpire WebSocket")

	// Subscribe to deposits channel
	subscribeMessage := []byte(`{
		"event": "subscribe",
		"data": {
			"channel": "deposits"
		}
	}`)

	if err := conn.WriteMessage(websocket.TextMessage, subscribeMessage); err != nil {
		log.Fatal("Subscribe error:", err)
	}

	// Message handling loop
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			log.Printf("Received update: %s\n", message)
		}
	}()

	// Connection maintenance
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("Ping error:", err)
				return
			}
		case <-interrupt:
			log.Println("Closing connection...")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Close error:", err)
			}
			return
		}
	}
}