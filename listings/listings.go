package listings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"  
	"strings"
)

type ListedItem struct {
	ID                 int64          `json:"id"`
	MarketName         string         `json:"market_name"`
	MarketValue        float64        `json:"market_value"`
	SuggestedPrice     float64        `json:"suggested_price"`
	Price              float64        `json:"purchase_price"`  // Renamed from Price
	Wear               float64        `json:"wear"`
	PublishedAt        string         `json:"published_at"`    // Was CreatedAt
	AboveRecommended   float64        `json:"above_recommended_price"`
	DepositorStats     DepositorStats `json:"depositor_stats"`
	ImageURL           string         `json:"preview_id"`      // Might need URL formatting
	WearName           string         `json:"wear_name"`
	Stickers           []interface{}  `json:"stickers"`        // Use specific struct if needed
	NameColor          string         `json:"name_color"`
	// Add other fields as needed
}

type DepositorStats struct {
	DeliveryRateRecent      float64 `json:"delivery_rate_recent"`
	DeliveryRateLong        float64 `json:"delivery_rate_long"`
	DeliveryTimeRecent      int     `json:"delivery_time_minutes_recent"`
	DeliveryTimeLong        int     `json:"delivery_time_minutes_long"`
	SteamLevelMin           int     `json:"steam_level_min_range"`
	SteamLevelMax           int     `json:"steam_level_max_range"`
	TradeNotifications      bool    `json:"user_has_trade_notifications_enabled"`
	OnlineStatus            int     `json:"user_online_status"`
}

type Pagination struct {
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	FirstPageURL string `json:"first_page_url"`
	LastPageURL  string `json:"last_page_url"`
	From        int    `json:"from"`
	// Add other pagination fields as needed
}


type APIResponse struct {
	Data       []ListedItem `json:"data"`
	Pagination Pagination   `json:"meta"`
}

// GetListedItems fetches paginated listings
func GetListedItems(page int) ([]ListedItem, *Pagination, error) {
	apiKey := os.Getenv("CSGO_EMPIRE_API_KEY")
	if apiKey == "" {
		return nil, nil, fmt.Errorf("missing API key")
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://csgoempire.com/api/v2/trading/items?per_page=100&state=2&page=%d",
			page,
		),
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("request creation failed: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var response APIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, nil, fmt.Errorf("JSON parse error: %v", err)
	}

	fmt.Println("Raw API response:", string(body))

	return response.Data, &response.Pagination, nil
}

func FilterItems(items []ListedItem) []ListedItem {
	var filteredItems []ListedItem
	for _, item := range items {
		if strings.Contains(item.MarketName, "StatTrak™") ||
			strings.Contains(item.MarketName, "Sticker") ||
			strings.Contains(item.MarketName, "Tag") ||
			strings.Contains(item.MarketName, "Souvenir") {
			continue // Skip items containing the specified substrings
		}
		filteredItems = append(filteredItems, item)
	}
	return filteredItems
}