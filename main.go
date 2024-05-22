// USAGE EXAMPLE

package main

import (
	"log"
	"net/http"
	"time"
)
import "github.com/mgrigoriev/goauth/authclient"

const token = "valid-token"
const authURL = "http://localhost:8083/api/v1/users/auth"
const timeout = 5 * time.Second

func main() {
	httpClient := &http.Client{Timeout: timeout}
	cfg := authclient.Config{AuthURL: authURL}
	cl := authclient.New(cfg, httpClient)

	userID, err := cl.Authenticate(token)
	if err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	log.Printf("authenticated user ID: %d", userID)
}
