package authclient

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
	delay      time.Duration
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	// Read and store the original request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	for attempts := 1; attempts <= rt.maxRetries; attempts++ {
		if bodyBytes != nil {
			// Reset the request body for each attempt
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		res, err = rt.next.RoundTrip(req)

		if err == nil && res.StatusCode < http.StatusInternalServerError {
			log.Printf("status: %v", res.StatusCode)
			break
		}

		log.Printf("attempt %d, err: %v", attempts, err)

		select {
		case <-req.Context().Done():
			return res, req.Context().Err()
		case <-time.After(rt.delay):
		}
	}

	return
}
