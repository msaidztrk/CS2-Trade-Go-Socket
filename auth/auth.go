package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func Authenticate() (string, error) {
	apiKey := os.Getenv("API_KEY")


	url := "https://csgoempire.com/api/v2/user/socket/auth"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	token, ok := result["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}