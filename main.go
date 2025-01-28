package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"CS2-Trade-Go-Socket/listings"
)

func coinFormat(v float64) string {
    return fmt.Sprintf("coin %.2f", v)
}

func wearFormat(name string, value float64) string {
    return fmt.Sprintf("%s (%.4f)", name, value)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	currentPage := 1
	totalItems := 0

	for {
		items, pagination, err := listings.GetListedItems(currentPage)
		if err != nil {
			log.Fatalf("Failed to fetch listings: %v", err)
		}

		if len(items) == 0 {
			break
		}

		// Filter out items containing "StatTrak™", "Sticker", "Tag", or "Souvenir"
		filteredItems := listings.FilterItems(items)
		totalItems += len(filteredItems)
		fmt.Printf("\n=== Page %d/%d (%d items) ===\n", 
			pagination.CurrentPage,
			pagination.LastPage,
			len(filteredItems),
		)

		for _, item := range filteredItems {
			println(
				item.ID , " | Name : " ,  truncateString(item.MarketName, 45) , 
				" Price: " , coinFormat(item.Price) , " | Market Value: " , coinFormat(item.MarketValue)  , " | Suggested : ", coinFormat(item.SuggestedPrice) , "\n"+
				"Profit: " , coinFormat(item.SuggestedPrice-item.Price) , " | Above Recommended: " ,  coinFormat(item.AboveRecommended) , "\n" +
				"Wear: ",wearFormat(item.WearName, item.Wear) ,"\n",
			)
		}

		if currentPage >= pagination.LastPage {
			break
		}
		currentPage++
	}

	fmt.Printf("\nTotal listed items found: %d\n", totalItems)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTime(isoTime string) string {
	t, err := time.Parse(time.RFC3339Nano, isoTime)
	if err != nil {
		return isoTime
	}
	return t.Format("2006-01-02 15:04")
}