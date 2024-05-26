package authclient

import (
	"context"
	"golang.org/x/time/rate"
	"net/http"
)

type rateLimitRoundTripper struct {
	next    http.RoundTripper
	limiter *rate.Limiter
}

func (rl *rateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := rl.limiter.Wait(context.Background()); err != nil {
		return nil, err
	}

	return rl.next.RoundTrip(req)
}
