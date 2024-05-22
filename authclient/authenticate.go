package authclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (ac *Client) Authenticate(token string) (user *CurrentUser, err error) {
	data := map[string]string{"token": token}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := ac.HTTPClient.Post(ac.Cfg.AuthURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
