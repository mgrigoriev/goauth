package authclient

import (
	"github.com/sony/gobreaker"
	"net/http"
)

type circuitBreakerRoundTripper struct {
	next           http.RoundTripper
	circuitBreaker *gobreaker.CircuitBreaker
}

func (cb *circuitBreakerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	_, err = cb.circuitBreaker.Execute(func() (interface{}, error) {
		resp, err = cb.next.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	})

	return resp, err
}
