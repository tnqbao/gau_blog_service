package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUserByID(userID uint64) (*User, error) {
	url := fmt.Sprintf("https://api.daudoo.com/api/user/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user: %s", resp.Status)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
