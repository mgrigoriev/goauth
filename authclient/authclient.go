package authclient

import (
	"net/http"
	"time"
)

const authURL = "http://users:8080/api/v1/users/auth"
const timeout = 5 * time.Second

type currentUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	httpClient := http.Client{Timeout: timeout}

	return &Client{
		httpClient: &httpClient,
	}
}
