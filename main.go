package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"CS2-Trade-Go-Socket/listings"
)

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

		totalItems += len(items)
		fmt.Printf("Page %d/%d (%d items)\n", 
			pagination.CurrentPage,
			pagination.LastPage,
			len(items),
		)

		for _, item := range items {
			fmt.Printf(
				"[%d] %-40s $%7.2f  %s\n",
				item.ID,
				truncateString(item.MarketName, 35),
				item.Price/100,
				formatTime(item.CreatedAt),
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
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return "invalid timestamp"
	}
	return t.Format("2006-01-02 15:04")
}