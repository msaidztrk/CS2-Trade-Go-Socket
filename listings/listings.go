package listings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ListedItem struct {
	ID         int     `json:"id"`
	MarketName string  `json:"market_name"`
	Price      float64 `json:"price"`
	Status     int     `json:"status"`
	CreatedAt  string  `json:"created_at"`
	ImageURL   string  `json:"image"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
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

	return response.Data, &response.Pagination, nil
}