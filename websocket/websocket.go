package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Connect to the WebSocket server
func Connect(url string, headers map[string]string) (*websocket.Conn, error) {
	wsHeaders := http.Header{}
	for key, value := range headers {
		wsHeaders.Set(key, value)
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, wsHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %v", err)
	}

	log.Println("WebSocket connection established.")
	return conn, nil
}

// Subscribe to a WebSocket channel
func Subscribe(conn *websocket.Conn, channel string) error {
	subscribeMessage := map[string]interface{}{
		"type":    "subscribe",
		"channel": channel,
	}
	if err := conn.WriteJSON(subscribeMessage); err != nil {
		return fmt.Errorf("failed to send subscribe message: %v", err)
	}

	log.Printf("Subscribed to channel: %s", channel)
	return nil
}

// Read messages from the WebSocket connection
func ReadMessages(conn *websocket.Conn, handler func(map[string]interface{})) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		handler(data)
	}
}