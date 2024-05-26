// USAGE EXAMPLE

package main

import (
	"log"
)
import "github.com/mgrigoriev/goauth/authclient"

const token = "valid-token"
const authURL = "http://localhost:8080/api/v1/users/auth"

func main() {
	cfg := authclient.Config{AuthURL: authURL}
	cl := authclient.New(cfg)

	// Uncomment "for {}" loop to test rate limiter
	//for {
	user, err := cl.Authenticate(token)
	if err != nil {
		log.Printf("authentication failed: %v", err)
	}

	if user != nil {
		log.Printf("authenticated user ID: %d", user.ID)
	}
	//}
}
