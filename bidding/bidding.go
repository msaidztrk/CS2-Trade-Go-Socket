package bidding

import (
	"log"

	"github.com/gorilla/websocket"
)

// Handle bidding updates
func HandleBiddingUpdate(conn *websocket.Conn, data map[string]interface{}) {
	biddingData := data["data"].(map[string]interface{})
	itemID := biddingData["item_id"].(string)
	price := biddingData["price"].(float64)

	log.Printf("Bidding update - Item ID: %s, Price: %.2f", itemID, price)

	// Example: Place a bid if the price is below a certain threshold
	if price < 100 { // Replace with your logic
		bidMessage := map[string]interface{}{
			"type": "bid",
			"data": map[string]interface{}{
				"item_id": itemID,
				"amount":  price + 1, // Bid 1 cent higher
			},
		}
		if err := conn.WriteJSON(bidMessage); err != nil {
			log.Printf("Failed to place bid: %v", err)
		} else {
			log.Printf("Placed bid on item %s for %.2f cents.", itemID, price+1)
		}
	}
}